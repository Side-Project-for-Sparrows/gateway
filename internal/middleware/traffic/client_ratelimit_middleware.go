package traffic

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/traffic/slidingwindow"
)

var ClientLimiter Limiter

type ClientRateLimiterInitializer struct{}

func (c *ClientRateLimiterInitializer) Construct() error {
	log.Println("[Construct] ClientRateLimiter Initialize 호출")

	ClientLimiter = slidingwindow.NewRateLimiter()
	return nil
}

func init() {
	config.Register(&ClientRateLimiterInitializer{})
}

func ClientRateLimitMiddleware() middlewaretype.Middleware {
	return func(input middlewaretype.MiddlewareInput) (*middlewaretype.HeaderPatch, error) {
		ip := extractIP(input)
		log.Printf("client rate limit")
		if ClientLimiter.IsOverLimit(ip, time.Now()) {
			log.Printf("[RateLimit] ip=%s blocked", ip)
			return &middlewaretype.HeaderPatch{}, fmt.Errorf("rate limit exceeded of client")
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
