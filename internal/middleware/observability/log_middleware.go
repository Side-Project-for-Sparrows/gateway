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
	return func(w http.ResponseWriter, r *http.Request) error {
		start := time.Now()
		defer func() {
			latency := time.Since(start)
			w.Header().Set(TimeHeader, latency.String())
			log.Printf("[%s] %s %s (%s)", r.Method, r.URL.Path, r.Header.Get("X_TRACE_ID"), latency)
		}()
		return nil
	}
}
