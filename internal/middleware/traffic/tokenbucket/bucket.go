package tokenbucket

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config/ratelimit"
)

type Bucket struct {
	Tokens     atomic.Int64
	Capacity   int64
	RefillRate int64
	LastUpdate atomic.Pointer[time.Time]
	Mu         sync.Mutex // 필요시 음수 방지 용도
	refilling  atomic.Bool
}

func NewBucket() *Bucket {
	capacity := ratelimit.Config.TokenBucket.Capacity
	refillRate := ratelimit.Config.TokenBucket.RefillRate

	b := &Bucket{
		Capacity:   capacity,
		RefillRate: refillRate,
	}
	b.Tokens.Store(capacity)
	now := time.Now()
	b.LastUpdate.Store(&now)
	return b
}

func (b *Bucket) TryRequest() bool {
	now := time.Now()

	// 1초 이상 지났고, 리필 중이 아니라면 고루틴으로 리필
	last := b.LastUpdate.Load()
	if last != nil && now.Sub(*last) >= time.Second && b.refilling.CompareAndSwap(false, true) {
		go b.refill(now)
	}

	new := b.Tokens.Add(-1)
	if new >= 0 {
		return true
	}

	b.Tokens.Add(1)
	return false
}

func (b *Bucket) refill(now time.Time) {
	defer b.refilling.Store(false)

	last := b.LastUpdate.Load()
	if last == nil {
		b.LastUpdate.Store(&now)
		return
	}

	elapsed := now.Sub(*last).Seconds()
	if elapsed <= 0 {
		return
	}

	newTokens := int64(elapsed * float64(b.RefillRate))
	if newTokens <= 0 {
		return
	}

	curr := b.Tokens.Load()
	refilled := min(curr+newTokens, b.Capacity)
	b.Tokens.Store(refilled)

	b.LastUpdate.Store(&now)
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
