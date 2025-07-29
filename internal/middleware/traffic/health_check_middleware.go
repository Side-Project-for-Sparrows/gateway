package traffic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
)

// 서비스 이름은 "user-service" / "board-service" 같은 형식이라고 가정
func HealthCheckMiddleware(hs *HealthStatus) middlewaretype.Middleware {
	return func(input middlewaretype.MiddlewareInput) (*middlewaretype.HeaderPatch, error) {
		serviceName := extractServiceFromPath(input.Path())
		if serviceName == "" {
			// 어떤 서비스로 가는 요청인지 식별 불가 → 통과시킴
			return nil, nil
		}

		// 서비스가 unhealthy인 경우 요청 차단
		if !hs.IsHealthy(serviceName) {
			body, _ := json.Marshal(map[string]string{
				"reason": fmt.Sprintf("service %s is currently unhealthy", serviceName),
			})

			return &middlewaretype.HeaderPatch{
				ResponseAdd:        http.Header{"Content-Type": []string{"application/json"}},
				ResponseStatusCode: http.StatusServiceUnavailable,
				ResponseBody:       body,
			}, fmt.Errorf("service %s is unhealthy", serviceName)
		}

		// 통과
		return nil, nil
	}
}
