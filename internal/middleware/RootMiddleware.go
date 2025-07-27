package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/observability"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/security"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/traffic"
)

func RootMiddlewareHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := middlewaretype.GenerateMiddlewareInput(r)
		// 체인 생성
		observeChain := middlewaretype.NewMiddlewareChain().
			AndThen(observability.TIDMiddleware()).
			AndThen(observability.LogMiddleware())

		parallel := middlewaretype.NewParallelChains().
			AndThen(security.JWTAuthMiddleware()).
			AndThen(traffic.RateLimitMiddleware())

		// 직렬 + 병렬을 AsMiddleware()로 묶고, 마지막 직렬로 구성
		fullChain := middlewaretype.NewMiddlewareChain().
			AndThen(parallel.AsMiddleware()).
			AndThen(observeChain.AsMiddleware())

		// Execute
		patches, err := fullChain.Execute(input)
		if err != nil {
			log.Printf("[RootMiddleware] execution failed: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Patch 적용
		ApplyPatches(r, w, patches)

		// 실제 핸들러 실행
		next.ServeHTTP(w, r)
	})
}

func ApplyPatches(r *http.Request, w http.ResponseWriter, patches []*middlewaretype.HeaderPatch) {
	for _, p := range patches {
		// 요청 헤더 추가
		for key, values := range p.RequestAdd {
			for _, v := range values {
				r.Header.Add(key, v)
			}
		}
		// 요청 헤더 삭제
		for _, key := range p.RequestDelete {
			r.Header.Del(key)
		}
		// 응답 헤더 추가
		for key, values := range p.ResponseAdd {
			for _, v := range values {
				w.Header().Add(key, v)
			}
		}
		// 응답 헤더 삭제
		for _, key := range p.ResponseDelete {
			w.Header().Del(key)
		}
		// context 값 추가
		if len(p.ContextAdd) > 0 {
			ctx := r.Context()
			for k, v := range p.ContextAdd {
				ctx = context.WithValue(ctx, k, v)
			}
			*r = *r.WithContext(ctx)
		}
	}
}
