package traffic

import (
	"log"
	"sync"
	"time"
)

type Window struct {
	Time  time.Time
	Count int
	Mu    sync.Mutex
}

func (w *Window) Increment() *Window {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	w.Count++
	return w
}

func newWindow(t time.Time) *Window {
	return &Window{
		Time:  t,
		Count: 0,
		Mu:    sync.Mutex{},
	}
}

type Windows struct {
	Prev *Window
	Curr *Window
	Mu   sync.Mutex
}

func (ws *Windows) Refresh(now time.Time) *Windows {
	ws.Mu.Lock()
	defer ws.Mu.Unlock()

	if now.After(ws.Curr.Time.Add(2 * time.Second)) {
		log.Printf("total refresh")
		return newWindows(now) // 완전 초기화
	}

	if now.After(ws.Curr.Time.Add(1 * time.Second)) {
		log.Printf("shitf")
		return shift(ws.Curr, now) // Curr → Prev, 새 Curr
	}

	ws.Curr.Increment()

	log.Printf("remain")

	return ws
}

func (ws *Windows) RateAt(t time.Time) int {
	log.Printf(t.String())
	now := t.Truncate(time.Millisecond) // 밀리초 정밀도 유지

	currStart := ws.Curr.Time

	// 예외 처리: Curr가 t보다 미래일 경우 → 논리 위반
	if now.Before(currStart) {
		return 0 // or panic("Curr time is in the future")
	}

	elapsed := now.Sub(currStart).Seconds() // 0.0 ~ 1.0

	if elapsed >= 1.0 {
		// 1초 이상 지났으면 Curr만 고려 (슬라이딩 안 함)
		ws.Curr.Mu.Lock()
		count := ws.Curr.Count
		ws.Curr.Mu.Unlock()
		return count
	}

	ws.Curr.Mu.Lock()
	currCount := ws.Curr.Count
	ws.Curr.Mu.Unlock()

	ws.Prev.Mu.Lock()
	prevCount := ws.Prev.Count
	ws.Prev.Mu.Unlock()

	rate := float64(currCount)*elapsed + float64(prevCount)*(1.0-elapsed)
	return int(rate + 0.5) // 반올림
}

func newWindows(t time.Time) *Windows {
	now := t.Truncate(time.Second)
	return &Windows{
		Prev: newWindow(now.Add(-1 * time.Second)),
		Curr: newWindow(now),
	}
}

func shift(w *Window, t time.Time) *Windows {
	now := t.Truncate(time.Second)
	return &Windows{
		Prev: w,
		Curr: newWindow(now),
	}
}
