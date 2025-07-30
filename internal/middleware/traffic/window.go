package traffic

import (
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
