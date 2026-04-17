// Package pause provides a pausable gate that can temporarily halt worker
// goroutines without cancelling the parent context.
package pause

import "sync"

// Controller allows callers to pause and resume a set of workers.
type Controller struct {
	mu     sync.Mutex
	cond   *sync.Cond
	paused bool
}

// New returns a ready-to-use Controller (initially resumed).
func New() *Controller {
	c := &Controller{}
	c.cond = sync.NewCond(&c.mu)
	return c
}

// Pause causes subsequent calls to Wait to block until Resume is called.
func (c *Controller) Pause() {
	c.mu.Lock()
	c.paused = true
	c.mu.Unlock()
}

// Resume unblocks all goroutines currently waiting in Wait.
func (c *Controller) Resume() {
	c.mu.Lock()
	c.paused = false
	c.cond.Broadcast()
	c.mu.Unlock()
}

// IsPaused reports whether the controller is currently paused.
func (c *Controller) IsPaused() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.paused
}

// Wait blocks the caller while the controller is paused.
// It returns immediately when the controller is resumed.
func (c *Controller) Wait() {
	c.mu.Lock()
	for c.paused {
		c.cond.Wait()
	}
	c.mu.Unlock()
}
