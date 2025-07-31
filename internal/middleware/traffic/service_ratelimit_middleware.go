package traffic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/traffic/tokenbucket"
	"github.com/Side-Project-for-Sparrows/gateway/internal/util"
)

var ServiceLimiter Limiter

func init() {
	ServiceLimiter = tokenbucket.NewRateLimiter(10 * time.Second)
}

func ServiceRateLimitMiddleware() middlewaretype.Middleware {
	return func(input middlewaretype.MiddlewareInput) (*middlewaretype.HeaderPatch, error) {
		serviceName := util.ExtractServiceKey(input.Path())

		if ServiceLimiter.IsOverLimit(serviceName, time.Now()) {
			log.Printf("[RateLimit] service=%s blocked", serviceName)
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
