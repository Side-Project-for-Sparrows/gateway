package circuitbreaker

import (
	"log"
	"net"

	"github.com/Side-Project-for-Sparrows/gateway/config/circuitbreak"
)

type circuitbreaker struct {
	circuits map[string]*circuit
}

func (c *circuitbreaker) IsHealthy(svc string) bool {
	circuit := c.circuits[svc]
	return circuit.IsHealthy()
}

func (c *circuitbreaker) ping(svc string) {
	conn, err := net.DialTimeout("tcp", svc, circuitbreak.Config.PingInterval)
	if err == nil {
		conn.Close()
	}

	healthy := err == nil
	circuit := c.circuits[svc].set
	hs.Set(svc+"-service", healthy)
	if !healthy {
		log.Printf("[HealthCheck] %s is UNHEALTHY (err=%v)", svc, err)
	} else {
		log.Printf("[HealthCheck] %s is healthy", svc)
	}
}
