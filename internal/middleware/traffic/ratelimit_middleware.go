package traffic

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config/ratelimit"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
)

var (
	rateLimitMap = make(map[string]*Windows)
)

func RateLimitMiddleware() middlewaretype.Middleware {
	return func(w http.ResponseWriter, r *http.Request) error {
		ip := getIpFrom(r)
		if isOverRateLimit(ip) {
			respMap := map[string]any{
				"reason": "TOO MANY REQUEST",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(respMap)
			return fmt.Errorf("rate limit exceeded")
		}

		//log.Printf("uner limit")
		return nil
	}
}

func getIpFrom(r *http.Request) string {
	xff := r.Header.Get("X-Forwared-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
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
	log.Printf("[RateLimit] IP=%s count=%d", ip, rate)

	return rate >= ratelimit.Config.Limit
}
