// Package adaptive provides a concurrency controller that adjusts the number
// of workers up or down based on observed error rate.
package adaptive

import (
	"sync"
	"sync/atomic"
)

// Controller adjusts a concurrency level based on success/failure feedback.
type Controller struct {
	mu      sync.Mutex
	current int64
	min     int64
	max     int64
	step    int64

	total  atomic.Int64
	errors atomic.Int64

	// ErrorThreshold is the ratio above which concurrency is reduced.
	ErrorThreshold float64
}

// New returns a Controller with the given bounds and step size.
// ErrorThreshold defaults to 0.1 (10%).
func New(min, max, step int) *Controller {
	if min < 1 {
		min = 1
	}
	if max < min {
		max = min
	}
	if step < 1 {
		step = 1
	}
	return &Controller{
		current:        int64(min),
		min:            int64(min),
		max:            int64(max),
		step:           int64(step),
		ErrorThreshold: 0.1,
	}
}

// Record registers the outcome of a single request.
func (c *Controller) Record(err bool) {
	c.total.Add(1)
	if err {
		c.errors.Add(1)
	}
}

// Adjust evaluates the current error rate and increases or decreases
// concurrency accordingly. It resets the counters after each call.
func (c *Controller) Adjust() int {
	total := c.total.Swap(0)
	errors := c.errors.Swap(0)

	c.mu.Lock()
	defer c.mu.Unlock()

	if total == 0 {
		return int(c.current)
	}

	rate := float64(errors) / float64(total)
	if rate > c.ErrorThreshold {
		c.current -= c.step
		if c.current < c.min {
			c.current = c.min
		}
	} else {
		c.current += c.step
		if c.current > c.max {
			c.current = c.max
		}
	}
	return int(c.current)
}

// Current returns the current concurrency level without adjusting.
func (c *Controller) Current() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return int(c.current)
}
