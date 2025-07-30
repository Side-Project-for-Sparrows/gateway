package traffic

import (
	"fmt"
	"log"
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
			for name, baseURL := range route.RouteMap {
				healthURL := fmt.Sprintf("%s/healthz", baseURL)

				log.Printf("HEL URL : %s", healthURL)
				go func(svc string, url string) {
					client := &http.Client{Timeout: 1 * time.Second}
					resp, err := client.Get(url)
					healthy := err == nil && resp.StatusCode == http.StatusOK
					hs.Set(svc+"-service", healthy)
					if !healthy {
						log.Printf("[HealthCheck] %s is UNHEALTHY (err=%v)", svc, err)
					} else {
						log.Printf("[HealthCheck] %s is healthy", svc)
					}

				}(name, healthURL)
			}
		}
	}()
}
