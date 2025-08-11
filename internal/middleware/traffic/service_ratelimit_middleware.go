package traffic

import (
	"fmt"
	"log"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/traffic/tokenbucket"
	"github.com/Side-Project-for-Sparrows/gateway/internal/util"
)

var ServiceLimiter Limiter

type ServiceRateLimiterInitializer struct{}

func (c *ServiceRateLimiterInitializer) Construct() error {
	log.Println("[Construct] ServiceRateLimiter Initialize 호출")

	ServiceLimiter = tokenbucket.NewRateLimiter()
	return nil
}

func init() {
	config.Register(&ServiceRateLimiterInitializer{})
}

func ServiceRateLimitMiddleware() middlewaretype.Middleware {
	return func(input middlewaretype.MiddlewareInput) (*middlewaretype.HeaderPatch, error) {
		serviceName := util.ExtractServiceKey(input.Path())

		if ServiceLimiter.IsOverLimit(serviceName, time.Now()) {
			log.Printf("[RateLimit] service=%s blocked", serviceName)
			return &middlewaretype.HeaderPatch{}, fmt.Errorf("rate limit exceeded on service ")
		}
		return &middlewaretype.HeaderPatch{}, nil
	}
}
