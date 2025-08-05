package middlewaretype

import (
	"sync"
)

type ParallelMiddleware struct {
	middlewares []Middleware
}

func NewParallelChains() *ParallelMiddleware {
	return &ParallelMiddleware{middlewares: []Middleware{}}
}

func (pc *ParallelMiddleware) AndThen(mw Middleware) *ParallelMiddleware {
	pc.middlewares = append(pc.middlewares, mw)
	return pc
}

// 락없는 구현
func (pc *ParallelMiddleware) Execute(input MiddlewareInput) ([]*HeaderPatch, error) {
	patches := make([]*HeaderPatch, len(pc.middlewares))
	errOnce := sync.Once{}
	var errRet error
	wg := sync.WaitGroup{}

	for i, mw := range pc.middlewares {
		wg.Add(1)

		go func(idx int, m Middleware) {
			defer wg.Done()

			patch, err := m(input)
			if err != nil {
				errOnce.Do(func() {
					errRet = err
				})
				return
			}
			patches[idx] = patch // mutex등 없이 race contition 방지하기 위해 배열로 저장
		}(i, mw)
	}

	wg.Wait()

	if errRet != nil {
		return nil, errRet
	}

	// 필터링 후 병합
	final := make([]*HeaderPatch, 0, len(patches))
	for _, p := range patches {
		if p != nil {
			final = append(final, p)
		}
	}
	return final, nil
}

// 락있는 구현
func (pc *ParallelMiddleware) Execute1(input MiddlewareInput) ([]*HeaderPatch, error) {

	var (
		mu      sync.Mutex
		patches []*HeaderPatch
		errOnce sync.Once
		errRet  error
		wg      sync.WaitGroup
	)

	for _, mw := range pc.middlewares {
		wg.Add(1)
		go func(m Middleware) {
			defer wg.Done()

			patch, err := m(input)
			if err != nil {
				errOnce.Do(func() {
					errRet = err
				})
				return
			}
			if patch != nil {
				mu.Lock()
				patches = append(patches, patch)
				mu.Unlock()
			}
		}(mw)
	}

	wg.Wait()
	if errRet != nil {
		return nil, errRet
	}
	return patches, nil
}

func (pc *ParallelMiddleware) AsMiddleware() Middleware {
	return func(input MiddlewareInput) (*HeaderPatch, error) {
		patches, err := pc.Execute(input)
		if err != nil {
			return nil, err
		}
		return mergePatches(patches), nil
	}
}
