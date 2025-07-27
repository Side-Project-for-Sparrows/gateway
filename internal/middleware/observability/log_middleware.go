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

		path := input.Path()
		method := input.Method()

		latency := time.Since(start)

		log.Printf("request latency: [%s] %s (%s)", method, path, latency)

		return &middlewaretype.HeaderPatch{
			ResponseAdd: http.Header{
				TimeHeader: []string{latency.String()},
			},
		}, nil
	}
}
