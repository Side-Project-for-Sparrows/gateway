package middleware

import (
	"net/http"

	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/observability"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/security"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/traffic"
)

// RootMiddlewareHandler 는 모든 미들웨어 구조를 묶어 api.Use() 에 등록할 단 하나의 미들웨어
func RootMiddlewareHandler(next http.Handler) http.Handler {
	// 직렬 체인
	observeChain := middlewaretype.NewMiddlewareChain().
		AndThen(observability.TIDMiddleware()).
		AndThen(observability.LogMiddleware())

	// 병렬 그룹 구성
	parallel := middlewaretype.NewParallelChains().
		AndThen(security.JWTAuthMiddleware()).
		AndThen(traffic.RateLimitMiddleware())

	last := middlewaretype.NewMiddlewareChain().
		AndThen(parallel.AsMiddleware()).
		AndThen(observeChain.AsMiddleware())

	// 병렬 미들웨어 실행 후 다음 핸들러 실행
	return last.AsHandler(next)
}
