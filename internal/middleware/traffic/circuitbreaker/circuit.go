package circuitbreaker

import "sync"

type circuit struct {
	name   string
	status CircuitStatus
	mu     sync.RWMutex
}

func NewCircuit(name string) *circuit {
	return &circuit{name: name, status: CLOSED}
}

func (c *circuit) SetStatus(next CircuitStatus) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch c.status {
	case CLOSED:
		if next == OPEN {
			c.status = next
			return true
		}
	case OPEN:
		if next == HALF_OPEN {
			c.status = next
			return true
		}
	case HALF_OPEN:
		if next == OPEN || next == CLOSED {
			c.status = next
			return true
		}
	}
	return false
}

func (c *circuit) Status() CircuitStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}

func (c *circuit) IsHealthy() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status == CLOSED
}
