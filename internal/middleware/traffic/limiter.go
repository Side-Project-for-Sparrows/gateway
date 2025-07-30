package traffic

import (
	"log"
	"sync"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config/ratelimit"
)

type RateLimiter struct {
	mu       sync.Mutex
	windows  map[string]*Windows
	interval time.Duration
}

func NewRateLimiter(cleanupInterval time.Duration) *RateLimiter {
	s := &RateLimiter{
		windows:  make(map[string]*Windows),
		interval: cleanupInterval,
	}
	go s.startCleanupLoop()
	return s
}

func (s *RateLimiter) IsOverLimit(target string, now time.Time) bool {
	s.mu.Lock()
	ws, exists := s.windows[target]
	if !exists {
		log.Printf("not exist!")
		ws = newWindows(now)
		s.windows[target] = ws
	}
	s.mu.Unlock()

	s.windows[target] = ws.Refresh()

	return s.windows[target].IsOverRateLimit()
}

func (s *RateLimiter) startCleanupLoop() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanupOldEntries()
	}
}

func (s *RateLimiter) cleanupOldEntries() {
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	for key, win := range s.windows {
		if now.Sub(win.Curr.Time) > 2*ratelimit.Config.WindowSize {
			log.Printf("delete")
			delete(s.windows, key)
		}
	}
}
