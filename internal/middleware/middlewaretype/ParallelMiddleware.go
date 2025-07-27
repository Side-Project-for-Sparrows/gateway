package middlewaretype

import (
	"context"
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

// 병렬로 미들웨어 실행, patch 수집, 에러 처리
func (pc *ParallelMiddleware) Execute(input MiddlewareInput) ([]*HeaderPatch, error) {
	_, cancel := context.WithCancel(input.Ctx())
	defer cancel()

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
					cancel()
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
