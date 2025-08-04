package slidingwindow

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config/ratelimit"
)

type Windows struct {
	Curr       atomic.Pointer[Window]
	Prev       atomic.Pointer[Window]
	lastAccess atomic.Int64 // UnixNano
	Mu         sync.Mutex
}

func (ws *Windows) Refresh() *Windows {
	ws.Mu.Lock()
	defer ws.Mu.Unlock()

	now := time.Now()
	curr := ws.Curr.Load()

	// 예외 처리
	if curr == nil {
		return newWindows(now)
	}

	if now.After(curr.Time.Add(2 * ratelimit.Config.SlidingWindow.WindowSize)) {
		return newWindows(now)
	}

	if now.After(curr.Time.Add(ratelimit.Config.SlidingWindow.WindowSize)) {
		return shift(curr, now)
	}

	curr.Increment()

	ws.lastAccess.Store(time.Now().UnixNano())
	return ws
}

func (ws *Windows) IsOverRateLimit() bool {
	t := time.Now()

	return ws.requestAt(t) >= int(ratelimit.Config.SlidingWindow.RequestsPerSecond*ratelimit.Config.SlidingWindow.WindowSize.Seconds())
}

func (ws *Windows) requestAt(t time.Time) int {
	curr := ws.Curr.Load()
	prev := ws.Prev.Load()

	if curr == nil || prev == nil {
		return 0
	}

	currStart := curr.Time
	now := t

	if now.Before(currStart) {
		return 0
	}

	elapsed := now.Sub(currStart).Seconds() / ratelimit.Config.SlidingWindow.WindowSize.Seconds()

	if elapsed >= 1.0 {
		return int(curr.Count.Load())
	}

	currCount := float64(curr.Count.Load())
	prevCount := float64(prev.Count.Load())

	rate := currCount + prevCount*(1.0-elapsed)
	return int(rate + 0.5) // 반올림
}

func newWindows(t time.Time) *Windows {
	now := truncate(t)

	ws := &Windows{}
	ws.Prev.Store(newWindow(now.Add(-ratelimit.Config.SlidingWindow.WindowSize)))
	ws.Curr.Store(newWindow(now))

	ws.lastAccess.Store(time.Now().UnixNano())
	return ws
}

func shift(prev *Window, t time.Time) *Windows {
	now := truncate(t)

	ws := &Windows{}
	ws.Prev.Store(prev)
	ws.Curr.Store(newWindow(now))
	return ws
}

func truncate(time time.Time) time.Time {
	return time.Truncate(ratelimit.Config.SlidingWindow.WindowSize)
}

func (ws *Windows) isOld(now time.Time) bool {
	return now.Sub(time.Unix(0, ws.lastAccess.Load())) > 2*ratelimit.Config.SlidingWindow.WindowSize
}
