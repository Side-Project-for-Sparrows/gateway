package traffic

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config/ratelimit"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
)

var (
	rateLimitMap = make(map[string]*Windows)
	mu           sync.Mutex
)

func RateLimitMiddleware() middlewaretype.Middleware {
	return func(input middlewaretype.MiddlewareInput) (*middlewaretype.HeaderPatch, error) {
		ip := extractIP(input)

		if isOverRateLimit(ip) {
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

		return nil, nil
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

func isOverRateLimit(ip string) bool {
	ws, exists := rateLimitMap[ip]
	t := time.Now()

	if !exists {
		ws = newWindows(t)
	}

	rateLimitMap[ip] = ws.Refresh(t)

	rate := ws.RateAt(t)
	log.Printf("[RateLimit] IP=%s count=%d limit=%d", ip, rate, ratelimit.Config.Limit)

	return rate >= ratelimit.Config.Limit
}

func init() {
	go startCleanupLoop()
}

func startCleanupLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		cleanupOldEntries()
	}
}

func cleanupOldEntries() {
	now := time.Now()

	mu.Lock()
	defer mu.Unlock()

	for key, win := range rateLimitMap {
		if now.Sub(win.Curr.Time) > 2*time.Second {
			delete(serviceRateLimitMap, key)
		}
	}
}
