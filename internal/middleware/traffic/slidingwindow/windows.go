package slidingwindow

import (
	"sync"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config/ratelimit"
)

type Windows struct {
	CurrCount  int
	PrevCount  int
	CurrTime   time.Time
	lastAccess int64 // UnixNano
	Mu         sync.Mutex
}

func (ws *Windows) refresh(now time.Time) *Windows {
	//log.Printf("before refresh curr is %d prev count is %d currcount is %d", ws.CurrTime.Second(), ws.PrevCount, ws.CurrCount)

	if now.After(ws.CurrTime.Add(2 * ratelimit.Config.SlidingWindow.WindowSize)) {
		//log.Printf("cur is total new")
		ws.newWindows(now)
	} else if now.After(ws.CurrTime.Add(ratelimit.Config.SlidingWindow.WindowSize)) {
		//log.Printf("cur is shift")
		ws.shift(now)
	}

	ws.CurrCount++
	//log.Printf("after refresh curr is %d prev count is %d currcount is %d", ws.CurrTime.Second(), ws.PrevCount, ws.CurrCount)
	ws.lastAccess = time.Now().UnixNano()
	return ws
}

func (ws *Windows) IsOverRateLimit(now time.Time) bool {
	ws.Mu.Lock()
	defer ws.Mu.Unlock()
	//log.Printf("ratelimit start at time %d", now.Second())
	ws.refresh(now)
	return ws.requestAt(now) >= ratelimit.Config.SlidingWindow.RequestsPerSecond*ratelimit.Config.SlidingWindow.WindowSize.Seconds()
}

func (ws *Windows) requestAt(t time.Time) float64 {
	now := t

	if now.Before(ws.CurrTime) {
		return 0
	}

	elapsed := now.Sub(ws.CurrTime).Seconds() / ratelimit.Config.SlidingWindow.WindowSize.Seconds()
	//log.Printf("elapsed is %.5f", elapsed)
	if elapsed >= 1.0 {
		return float64(ws.CurrCount)
	}

	currCount := float64(ws.CurrCount)
	prevCount := float64(ws.PrevCount)

	//log.Printf("curr is %d prev count is %d currcount is %d", ws.CurrTime.Second(), ws.PrevCount, ws.CurrCount)

	rate := currCount + prevCount*(1.0-elapsed)
	//log.Printf("rate is %.5f", rate)

	return rate
}

func newWindows(t time.Time) *Windows {
	ws := &Windows{}
	ws.newWindows(t)
	return ws
}

func (ws *Windows) newWindows(t time.Time) {
	now := truncate(t)

	ws.PrevCount = 0
	ws.CurrCount = 0
	ws.CurrTime = now
	ws.lastAccess = time.Now().UnixNano()
}

func (ws *Windows) shift(t time.Time) {
	now := truncate(t)

	ws.PrevCount = ws.CurrCount
	ws.CurrCount = 0
	ws.CurrTime = now
	ws.lastAccess = time.Now().UnixNano()
}

func truncate(time time.Time) time.Time {
	return time.Truncate(ratelimit.Config.SlidingWindow.WindowSize)
}

func (ws *Windows) isOld(now time.Time) bool {
	return now.Sub(time.Unix(0, ws.lastAccess)) > 2*ratelimit.Config.SlidingWindow.WindowSize
}
