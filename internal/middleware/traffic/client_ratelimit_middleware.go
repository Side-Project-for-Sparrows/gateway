package traffic

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/traffic/slidingwindow"
)

var ClientLimiter Limiter

func init() {
	ClientLimiter = slidingwindow.NewRateLimiter(10 * time.Second)
}

func ClientRateLimitMiddleware() middlewaretype.Middleware {
	return func(input middlewaretype.MiddlewareInput) (*middlewaretype.HeaderPatch, error) {
		ip := extractIP(input)

		if ClientLimiter.IsOverLimit(ip, time.Now()) {
			log.Printf("[RateLimit] ip=%s blocked", ip)
			resp := map[string]any{"reason": "TOO MANY REQUEST"}
			body, _ := json.Marshal(resp)
			return &middlewaretype.HeaderPatch{
				ResponseAdd: http.Header{
					"Content-Type": []string{"application/json"},
				},
				ResponseStatusCode: http.StatusTooManyRequests,
				ResponseBody:       body,
			}, fmt.Errorf("rate limit exceeded")
		}

		return &middlewaretype.HeaderPatch{}, nil
	}
}

func extractIP(input middlewaretype.MiddlewareInput) string {
	h := input.Headers()
	xff := h.Get("X-Forwarded-For")

	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	ip, _, err := net.SplitHostPort(input.RemoteAddr())
	if err != nil {
		return input.RemoteAddr()
	}
	return ip
}
