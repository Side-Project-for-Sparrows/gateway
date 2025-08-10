package circuitbreaker

import (
	"container/list"
	"math"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config/circuitbreak"
)

type Event struct {
	Timestamp int64
	Success   bool
}

type Circuit struct {
	url        string
	events     *list.List
	state      atomic.Value
	failCnt    atomic.Int32
	scsCnt     atomic.Int32
	gradeBits  atomic.Uint64
	inProgress atomic.Bool
}

func NewCircuit(url string) *Circuit {
	c := &Circuit{
		url:    strings.TrimPrefix(url, "http://"),
		events: list.New(),
	}

	c.state.Store(&CircuitState{stateType: Closed})
	c.failCnt.Store(0)
	c.scsCnt.Store(0)
	return c
}

func (c *Circuit) IsHealthy() bool {
	return c.state.Load().(*CircuitState).IsHealthy()
}

func (c *Circuit) Next() {
	if !c.inProgress.CompareAndSwap(false, true) {
		return
	}
	defer c.inProgress.Store(false)

	_, err := net.DialTimeout("tcp", c.url, circuitbreak.Config.PingInterval)
	success := err == nil
	now := time.Now().Unix()

	c.updateEventQ(success, now)

	status := c.state.Load().(*CircuitState)
	nextType := status.Next(c)
	c.state.Store(&CircuitState{stateType: nextType})
}

func (c *Circuit) updateEventQ(success bool, now int64) {
	c.pruneQ(now)
	c.add(success, now)
}

func (c *Circuit) pruneQ(now int64) {
	cutoff := now - int64(circuitbreak.Config.EventQueueSize.Seconds())
	for e := c.events.Front(); e != nil; {
		ev := e.Value.(Event)
		if ev.Timestamp >= cutoff {
			break
		}
		next := e.Next()
		if ev.Success {
			c.scsCnt.Add(-1)
		} else {
			c.failCnt.Add(-1)
		}

		c.events.Remove(e)
		e = next
	}
}

func (c *Circuit) add(success bool, now int64) {
	c.events.PushBack(Event{Timestamp: now, Success: success})
	grade := c.getGrade() * circuitbreak.Config.Weight
	if success {
		c.scsCnt.Add(1)
		grade += 1
	} else {
		c.failCnt.Add(1)
	}
	c.setGrade(grade)
}

func (c *Circuit) setGrade(val float64) {
	c.gradeBits.Store(math.Float64bits(val))
}

func (c *Circuit) getGrade() float64 {
	return math.Float64frombits(c.gradeBits.Load())
}
