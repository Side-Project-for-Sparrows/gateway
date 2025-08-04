package circuitbreaker

import (
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config/circuitbreak"
	"github.com/Side-Project-for-Sparrows/gateway/config/route"
)

type CircuitBreaker struct {
	circuits map[string]*Circuit
}

func NewCircuitBreaker() CircuitBreaker {
	cb := CircuitBreaker{
		circuits: make(map[string]*Circuit),
	}
	for name, addr := range route.RouteMap {
		cb.circuits[name] = NewCircuit(addr)
	}
	go cb.startPingLoop()
	return cb
}

func (cb *CircuitBreaker) IsHealthy(svc string) bool {
	circuit, ok := cb.circuits[svc]
	if !ok {
		return false
	}
	return circuit.IsHealthy()
}

func (cb *CircuitBreaker) startPingLoop() {
	ticker := time.NewTicker(circuitbreak.Config.PingInterval)
	defer ticker.Stop()

	for range ticker.C {
		for _, c := range cb.circuits {
			c.Next()
		}
	}
}
