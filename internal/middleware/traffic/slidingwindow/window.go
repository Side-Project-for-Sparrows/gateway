package slidingwindow

import (
	"sync/atomic"
	"time"
)

type Window struct {
	Time  time.Time
	Count atomic.Int32
}

func (w *Window) Increment() {
	w.Count.Add(1)
}

func newWindow(t time.Time) *Window {
	return &Window{
		Time: t,
	}
}
