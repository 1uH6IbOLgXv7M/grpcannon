package watchdog

import "sync/atomic"

// Counter is a simple thread-safe Source backed by atomic counters.
type Counter struct {
	total  atomic.Int64
	errors atomic.Int64
}

// RecordSuccess increments the total request count.
func (c *Counter) RecordSuccess() {
	c.total.Add(1)
}

// RecordError increments both the total and error counts.
func (c *Counter) RecordError() {
	c.total.Add(1)
	c.errors.Add(1)
}

// Total returns the total number of recorded requests.
func (c *Counter) Total() int64 {
	return c.total.Load()
}

// ErrorRate returns the fraction of requests that were errors in [0,1].
func (c *Counter) ErrorRate() float64 {
	t := c.total.Load()
	if t == 0 {
		return 0
	}
	return float64(c.errors.Load()) / float64(t)
}
