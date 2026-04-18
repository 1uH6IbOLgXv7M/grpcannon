// Package window provides a sliding-window counter for tracking
// request counts and error rates over a rolling time interval.
package window

import (
	"sync"
	"time"
)

// bucket holds counts for a single time slice.
type bucket struct {
	total  int64
	errors int64
}

// Window is a sliding-window counter divided into fixed-size slots.
type Window struct {
	mu       sync.Mutex
	slots    []bucket
	times    []time.Time
	size     int
	slotSize time.Duration
	now      func() time.Time
}

// New creates a Window with the given total duration and number of slots.
func New(duration time.Duration, slots int) *Window {
	if slots < 1 {
		slots = 1
	}
	return &Window{
		slots:    make([]bucket, slots),
		times:    make([]time.Time, slots),
		size:     slots,
		slotSize: duration / time.Duration(slots),
		now:      time.Now,
	}
}

func (w *Window) slotIndex(t time.Time) int {
	return int(t.UnixNano()/int64(w.slotSize)) % w.size
}

// Add records a request outcome in the current time slot.
func (w *Window) Add(isError bool) {
	now := w.now()
	idx := w.slotIndex(now)
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.times[idx].IsZero() || now.Sub(w.times[idx]) >= w.slotSize {
		w.slots[idx] = bucket{}
		w.times[idx] = now
	}
	w.slots[idx].total++
	if isError {
		w.slots[idx].errors++
	}
}

// Counts returns the aggregate total and error counts across all live slots.
func (w *Window) Counts() (total, errors int64) {
	now := w.now()
	w.mu.Lock()
	defer w.mu.Unlock()
	for i := 0; i < w.size; i++ {
		if w.times[i].IsZero() || now.Sub(w.times[i]) >= w.slotSize*time.Duration(w.size) {
			continue
		}
		total += w.slots[i].total
		errors += w.slots[i].errors
	}
	return
}

// ErrorRate returns the fraction of requests that were errors (0 if no requests).
func (w *Window) ErrorRate() float64 {
	t, e := w.Counts()
	if t == 0 {
		return 0
	}
	return float64(e) / float64(t)
}
