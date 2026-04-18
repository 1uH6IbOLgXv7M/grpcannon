// Package debounce provides a debouncer that delays execution until
// a quiet period has elapsed since the last call.
package debounce

import (
	"sync"
	"time"
)

// Debouncer delays fn until no further calls arrive within interval.
type Debouncer struct {
	mu       sync.Mutex
	interval time.Duration
	timer    *time.Timer
	fn       func()
}

// New returns a Debouncer that fires fn after interval of inactivity.
// If interval is <= 0 it defaults to 100ms.
func New(interval time.Duration, fn func()) *Debouncer {
	if interval <= 0 {
		interval = 100 * time.Millisecond
	}
	return &Debouncer{interval: interval, fn: fn}
}

// Call schedules fn, resetting the timer if already pending.
func (d *Debouncer) Call() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.interval, func() {
		d.fn()
	})
}

// Flush fires fn immediately if a call is pending and cancels the timer.
// Returns true if a pending call was flushed.
func (d *Debouncer) Flush() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer == nil {
		return false
	}
	stopped := d.timer.Stop()
	d.timer = nil
	if stopped {
		d.fn()
		return true
	}
	return false
}

// Stop cancels any pending call without firing fn.
func (d *Debouncer) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}
