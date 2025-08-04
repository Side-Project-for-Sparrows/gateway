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

		rateLimitMiddleware := middlewaretype.NewParallelChains().
			AndThen(traffic.ClientRateLimitMiddleware()).
			AndThen(traffic.ServiceRateLimitMiddleware())

		trafficMiddleware := middlewaretype.NewSerialMiddlewareChain().
			AndThen(traffic.CircuitBreakerMiddleware()).
			AndThen(rateLimitMiddleware.AsMiddleware())

		validateMiddleware := middlewaretype.NewParallelChains().
			AndThen(trafficMiddleware.AsMiddleware()).
			AndThen(security.JWTAuthMiddleware())

		fullChain := middlewaretype.NewSerialMiddlewareChain().
			AndThen(validateMiddleware.AsMiddleware()).
			AndThen(observability.TIDMiddleware())

		patches, err := fullChain.Execute(input)

		ApplyPatches(r, w, patches)
		if err != nil {
			log.Printf("[RootMiddleware] execution failed: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ParellelRootMiddlewareHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := middlewaretype.GenerateMiddlewareInput(r)

		//모든 미들웨어 병렬 구성
		fullChain := middlewaretype.NewParallelChains().
			AndThen(security.JWTAuthMiddleware()).
			AndThen(traffic.CircuitBreakerMiddleware()).
			AndThen(traffic.ClientRateLimitMiddleware()).
			AndThen(traffic.ServiceRateLimitMiddleware()).
			AndThen(observability.TIDMiddleware()).
			AndThen(observability.LogMiddleware())

		patches, err := fullChain.Execute(input)

		ApplyPatches(r, w, patches)
		if err != nil {
			log.Printf("[RootMiddleware] execution failed: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SerialRootMiddlewareHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := middlewaretype.GenerateMiddlewareInput(r)

		// 모든 미들웨어를 직렬로 구성
		fullChain := middlewaretype.NewSerialMiddlewareChain().
			AndThen(observability.TIDMiddleware()).
			AndThen(observability.LogMiddleware()).
			AndThen(security.JWTAuthMiddleware()).
			AndThen(traffic.ClientRateLimitMiddleware()).
			AndThen(traffic.ServiceRateLimitMiddleware())

		patches, err := fullChain.Execute(input)

		ApplyPatches(r, w, patches)
		if err != nil {
			log.Printf("[RootMiddleware] execution failed: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ApplyPatches(r *http.Request, w http.ResponseWriter, patches []*middlewaretype.HeaderPatch) {
	for _, p := range patches {
		for key, values := range p.RequestAdd {
			for _, v := range values {
				r.Header.Add(key, v)
			}
		}

		for _, key := range p.RequestDelete {
			r.Header.Del(key)
		}

		for key, values := range p.ResponseAdd {
			for _, v := range values {
				w.Header().Add(key, v)
			}
		}

		for _, key := range p.ResponseDelete {
			w.Header().Del(key)
		}

		if len(p.ContextAdd) > 0 {
			ctx := r.Context()
			for k, v := range p.ContextAdd {
				ctx = context.WithValue(ctx, k, v)
			}
			*r = *r.WithContext(ctx)
		}
	}
}
