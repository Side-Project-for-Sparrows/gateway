package traffic

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config/route"
)

type HealthStatus struct {
	mu     sync.RWMutex
	status map[string]bool // 서비스명 → 상태 (true = healthy)
}

func NewHealthStatus() *HealthStatus {
	return &HealthStatus{
		status: make(map[string]bool),
	}
}

func (h *HealthStatus) IsHealthy(service string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.status[service]
}

func (h *HealthStatus) Set(service string, healthy bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.status[service] = healthy
}

func StartHealthCheckLoop(hs *HealthStatus) {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// reflect 대신 명시적으로 순회
			services := map[string]string{
				"user":   route.Config.User,
				"board":  route.Config.Board,
				"school": route.Config.School,
				"search": route.Config.Search,
			}

			for name, baseURL := range services {
				healthURL := fmt.Sprintf("%s/healthz", baseURL)

				go func(svc string, url string) {
					client := &http.Client{Timeout: 1 * time.Second}
					resp, err := client.Get(url)
					healthy := err == nil && resp.StatusCode == http.StatusOK
					hs.Set(svc+"-service", healthy)
				}(name, healthURL)
			}
		}
	}()
}
