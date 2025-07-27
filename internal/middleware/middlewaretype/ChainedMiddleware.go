package middlewaretype

import (
	"log"
	"net/http"
)

// 미들웨어 체인 (순차적으로 실행되는 미들웨어 그룹)
type ChainedMiddleware struct {
	middlewares []Middleware
}

func NewMiddlewareChain() *ChainedMiddleware {
	return &ChainedMiddleware{middlewares: []Middleware{}}
}

func (mc *ChainedMiddleware) AndThen(mw Middleware) *ChainedMiddleware {
	mc.middlewares = append(mc.middlewares, mw)
	return mc
}

func (mc *ChainedMiddleware) Execute(w http.ResponseWriter, r *http.Request) error {
	for i, mw := range mc.middlewares {
		if err := mw(w, r); err != nil {
			log.Printf("[ChainedMiddleware] Middleware #%d failed: %v", i, err)
			return err
		}
	}
	return nil
}

// MiddlewareChain을 하나의 MiddlewareFunc처럼 변환
func (mc *ChainedMiddleware) AsMiddleware() Middleware {
	return func(w http.ResponseWriter, r *http.Request) error {
		return mc.Execute(w, r)
	}
}

func (mc *ChainedMiddleware) AsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := mc.Execute(w, r); err != nil {
			// error 처리 정책은 여기서 담당
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
