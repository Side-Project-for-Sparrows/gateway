package observability

import (
	"log"
	"net/http"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
)

const (
	TimeHeader = "X_TIME_COST"
)

func LogMiddleware() middlewaretype.Middleware {
	return func(input middlewaretype.MiddlewareInput) (*middlewaretype.HeaderPatch, error) {
		start := time.Now()

		// context에 trace ID가 있다고 가정
		traceID := input.Ctx().Value("X_TRACE_ID")
		path := "unknown"
		method := "UNKNOWN"

		// context에 요청 정보가 포함돼 있다면 추출
		if r, ok := input.Ctx().Value("REQUEST_META").(*http.Request); ok {
			method = r.Method
			path = r.URL.Path
		}

		latency := time.Since(start)

		log.Printf("[%s] %s %v (%s)", method, path, traceID, latency)

		return &middlewaretype.HeaderPatch{
			ResponseAdd: http.Header{
				TimeHeader: []string{latency.String()},
			},
		}, nil
	}
}
