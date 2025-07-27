package middlewaretype

import (
	"context"
	"net/http"
	"sync"
)

type ParallelMiddleware struct {
	middlewares []Middleware
	mu          sync.Mutex
}

func (pc *ParallelMiddleware) Execute(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	errCh := make(chan error, len(pc.middlewares))
	wg := sync.WaitGroup{}

	for _, mw := range pc.middlewares {
		wg.Add(1)
		go func(m Middleware) {
			defer wg.Done()
			if err := m(w, r.WithContext(ctx)); err != nil {
				errCh <- err
				cancel() // 실패 시 취소
			}
		}(mw)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

func (pc *ParallelMiddleware) AsMiddleware() Middleware {
	return func(w http.ResponseWriter, r *http.Request) error {
		return pc.Execute(w, r)
	}
}

func (pc *ParallelMiddleware) AsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := pc.Execute(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func NewParallelChains() *ParallelMiddleware {
	return &ParallelMiddleware{middlewares: []Middleware{}}
}

func (pc *ParallelMiddleware) AndThen(mw Middleware) *ParallelMiddleware {
	pc.middlewares = append(pc.middlewares, mw)
	return pc
}
