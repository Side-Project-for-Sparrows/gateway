package traffic

import (
	"log"
	"sync"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config/ratelimit"
)

type Windows struct {
	Prev *Window
	Curr *Window
	Mu   sync.Mutex
}

func (ws *Windows) Refresh() *Windows {
	ws.Mu.Lock()
	defer ws.Mu.Unlock()

	now := time.Now()

	if now.After(ws.Curr.Time.Add(2 * ratelimit.Config.WindowSize)) {
		return newWindows(now) // 완전 초기화
	}

	if now.After(ws.Curr.Time.Add(ratelimit.Config.WindowSize)) {
		return shift(ws.Curr, now) // Curr → Prev, 새 Curr
	}

	ws.Curr.Increment()
	return ws
}

func (ws *Windows) IsOverRateLimit() bool {
	t := time.Now()
	return ws.requestAt(t) >= ratelimit.Config.RequestsPerSecond*ratelimit.Config.WindowSize.Seconds()
}

func (ws *Windows) requestAt(t time.Time) float64 {
	now := t
	currStart := ws.Curr.Time

	// 예외 처리: Curr가 t보다 미래일 경우 → 논리 위반
	if now.Before(currStart) {
		return 0
	}

	elapsed := now.Sub(currStart).Seconds() / ratelimit.Config.WindowSize.Seconds()

	if elapsed >= 1.0 {
		ws.Curr.Mu.Lock()
		count := ws.Curr.Count
		ws.Curr.Mu.Unlock()
		return float64(count)
	}

	ws.Curr.Mu.Lock()
	currCount := ws.Curr.Count
	ws.Curr.Mu.Unlock()

	ws.Prev.Mu.Lock()
	prevCount := ws.Prev.Count
	ws.Prev.Mu.Unlock()

	rate := float64(currCount)*elapsed + float64(prevCount)*(1.0-elapsed)
	return rate + 0.5 // 반올림
}

func newWindows(t time.Time) *Windows {
	now := truncate(t)
	return &Windows{
		Prev: newWindow(now.Add(-ratelimit.Config.WindowSize)),
		Curr: newWindow(now),
	}
}

func shift(w *Window, t time.Time) *Windows {
	now := truncate(t)
	log.Printf("before truncate : %d", t.Nanosecond())
	log.Printf("truncate now : %d", now.Nanosecond())
	return &Windows{
		Prev: w,
		Curr: newWindow(now),
	}
}

func truncate(time time.Time) time.Time {
	return time.Truncate(ratelimit.Config.WindowSize)
}
