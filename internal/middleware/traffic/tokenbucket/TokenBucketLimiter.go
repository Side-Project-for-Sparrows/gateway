// tokenbucket/limiter.go
package tokenbucket

import (
	"log"
	"sync"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config/ratelimit"
)

type TokenBucketLimiter struct {
	buckets  sync.Map // key: string, value: *Bucket
	interval time.Duration
}

func NewRateLimiter() *TokenBucketLimiter {
	lim := &TokenBucketLimiter{
		interval: ratelimit.Config.TokenBucket.CleanInterval,
	}
	go lim.startCleanupLoop()
	return lim
}

func (l *TokenBucketLimiter) IsOverLimit(key string, now time.Time) bool {
	val, ok := l.buckets.Load(key)
	if !ok {
		newBucket := NewBucket()
		val, _ = l.buckets.LoadOrStore(key, newBucket)
	}
	b := val.(*Bucket)
	return !b.TryRequest() // 토큰 없으면 초과
}

func (l *TokenBucketLimiter) startCleanupLoop() {
	ticker := time.NewTicker(l.interval)
	defer ticker.Stop()

	for range ticker.C {
		l.cleanupOldEntries()
	}
}

func (l *TokenBucketLimiter) cleanupOldEntries() {
	now := time.Now()
	l.buckets.Range(func(key, value any) bool {
		b := value.(*Bucket)
		last := b.LastUpdate.Load()
		if last != nil && now.Sub(*last) > 2*time.Minute {
			l.buckets.Delete(key)
			log.Printf("deleted token bucket for key: %v", key)
		}
		return true
	})
}
