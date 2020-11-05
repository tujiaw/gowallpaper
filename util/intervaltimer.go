package util

import (
	"time"
)

type IntervalTimer struct {
	timer    *time.Timer
	d        time.Duration
	canceled chan struct{}
	fn       func()
}

func NewIntervalTimer(d time.Duration, fn func()) *IntervalTimer {
	return &IntervalTimer{d: d, fn: fn}
}

func (t *IntervalTimer) Start() {
	t.canceled = make(chan struct{})
	t.timer = time.NewTimer(t.d)
	go func(t *IntervalTimer) {
		for {
			select {
			case <-t.timer.C:
				t.fn()
				t.timer.Reset(t.d)
			case <-t.canceled:
				return
			}
		}
	}(t)
}

func (t *IntervalTimer) Stop() {
	if t.canceled != nil {
		close(t.canceled)
	}
	if t.timer != nil {
		t.timer.Stop()
	}
}
