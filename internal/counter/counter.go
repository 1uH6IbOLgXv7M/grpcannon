// Package counter provides a thread-safe request/error counter with
// atomic operations suitable for high-throughput load testing scenarios.
package counter

import (
	"sync/atomic"
)

// Counter tracks total requests and errors atomically.
type Counter struct {
	total  atomic.Int64
	errors atomic.Int64
}

// New returns a zeroed Counter.
func New() *Counter {
	return &Counter{}
}

// IncTotal increments the total request count by one.
func (c *Counter) IncTotal() {
	c.total.Add(1)
}

// IncErrors increments both the error count and the total count by one.
func (c *Counter) IncErrors() {
	c.errors.Add(1)
	c.total.Add(1)
}

// Total returns the current total request count.
func (c *Counter) Total() int64 {
	return c.total.Load()
}

// Errors returns the current error count.
func (c *Counter) Errors() int64 {
	return c.errors.Load()
}

// ErrorRate returns the fraction of requests that resulted in an error.
// Returns 0 if no requests have been recorded.
func (c *Counter) ErrorRate() float64 {
	t := c.total.Load()
	if t == 0 {
		return 0
	}
	return float64(c.errors.Load()) / float64(t)
}

// Reset atomically zeroes both counters.
func (c *Counter) Reset() {
	c.total.Store(0)
	c.errors.Store(0)
}

// Snapshot returns a point-in-time copy of the counter values.
type Snapshot struct {
	Total  int64
	Errors int64
}

// Snapshot captures the current counter state without resetting it.
func (c *Counter) Snapshot() Snapshot {
	return Snapshot{
		Total:  c.total.Load(),
		Errors: c.errors.Load(),
	}
}
