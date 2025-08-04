package traffic

import (
	"fmt"
	"log"

	"github.com/Side-Project-for-Sparrows/gateway/config"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/traffic/circuitbreaker"
	"github.com/Side-Project-for-Sparrows/gateway/internal/util"
)

var cb circuitbreaker.CircuitBreaker

type CircuitBreakerInitializer struct{}

func init() {
	config.Register(&CircuitBreakerInitializer{})
}

func (c *CircuitBreakerInitializer) Construct() error {
	log.Println("[Construct] ClientRateLimiter Initialize 호출")

	cb = circuitbreaker.NewCircuitBreaker()
	return nil
}

func CircuitBreakerMiddleware() middlewaretype.Middleware {
	return func(input middlewaretype.MiddlewareInput) (*middlewaretype.HeaderPatch, error) {
		serviceName := util.ExtractServiceKey(input.Path())
		if serviceName == "" {
			// 식별 불가 시 에러
			return &middlewaretype.HeaderPatch{}, fmt.Errorf("circuit breaker: service %s is not found", serviceName)
		}

		if !cb.IsHealthy(serviceName) {
			return &middlewaretype.HeaderPatch{}, fmt.Errorf("circuit breaker: service %s is unavailable", serviceName)
		}

		return nil, nil
	}
}
