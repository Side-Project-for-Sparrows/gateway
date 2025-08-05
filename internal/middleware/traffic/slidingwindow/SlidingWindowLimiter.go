package slidingwindow

import (
	"log"
	"sync"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config/ratelimit"
)

type SlidingWindowLimiter struct {
	windows  sync.Map // key: string, value: *entry
	interval time.Duration
}

func NewRateLimiter() *SlidingWindowLimiter {
	s := &SlidingWindowLimiter{
		interval: ratelimit.Config.SlidingWindow.CleanInterval,
	}
	go s.startCleanupLoop()
	return s
}

func (s *SlidingWindowLimiter) IsOverLimit(target string, now time.Time) bool {
	val, ok := s.windows.Load(target)
	if !ok {
		newEntry := newWindows(now)
		val, _ = s.windows.LoadOrStore(target, newEntry)
	}

	ws := val.(*Windows)
	ws.Refresh()
	return ws.IsOverRateLimit()
}

func (s *SlidingWindowLimiter) startCleanupLoop() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanupOldEntries()
	}
}

func (s *SlidingWindowLimiter) cleanupOldEntries() {
	now := time.Now()
	s.windows.Range(func(key, value any) bool {
		e := value.(*Windows)
		if e.isOld(now) {
			s.windows.Delete(key)
			log.Printf("deleted rate limiter for key: %v", key)
		}
		return true
	})
}
