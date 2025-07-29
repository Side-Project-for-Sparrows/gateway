package traffic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config/ratelimit"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
)

var (
	serviceRateLimitMap = make(map[string]*Windows)
	mu                  sync.Mutex
)

// 서비스별 처리율 제한 미들웨어
func RateLimitByServiceMiddleware() middlewaretype.Middleware {
	return func(input middlewaretype.MiddlewareInput) (*middlewaretype.HeaderPatch, error) {
		serviceName := extractServiceFromPath(input.Path())

		if serviceName == "" {
			log.Printf("[RateLimit] Unknown path=%s → 서비스 식별 불가, 우회 처리", input.Path())
			return nil, nil
		}

		if isOverServiceRateLimit(serviceName) {
			resp := map[string]any{"reason": fmt.Sprintf("TOO MANY REQUEST to service: %s", serviceName)}
			body, _ := json.Marshal(resp)

			return &middlewaretype.HeaderPatch{
				ResponseAdd: http.Header{
					"Content-Type": []string{"application/json"},
				},
				ResponseStatusCode: http.StatusTooManyRequests,
				ResponseBody:       body,
			}, fmt.Errorf("rate limit exceeded for service: %s", serviceName)
		}

		return nil, nil
	}
}

// 요청 path 기반으로 서비스명 추출
func extractServiceFromPath(path string) string {
	// ex) "/user/profile" → "user-service"
	segments := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(segments) == 0 {
		return ""
	}

	return segments[0]
}

func isOverServiceRateLimit(service string) bool {
	ws, exists := serviceRateLimitMap[service]
	t := time.Now()

	if !exists {
		ws = newWindows(t)
	}

	serviceRateLimitMap[service] = ws.Refresh(t)

	rate := ws.RateAt(t)
	log.Printf("[RateLimit] service=%s count=%d", service, rate)

	return rate >= ratelimit.Config.Limit
}

func init() {
	go startCleanupLoop()
}

func startCleanupLoopFromService() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		cleanupOldEntries()
	}
}

func cleanupOldEntriesFromSerivce() {
	now := time.Now()

	mu.Lock()
	defer mu.Unlock()

	for key, win := range serviceRateLimitMap {
		if now.Sub(win.curr) > 2*time.Second {
			delete(serviceRateLimitMap, key)
		}
	}
}
