// Package cooldown provides a time-based cooldown tracker that prevents
// actions from being taken more frequently than a configured interval.
package cooldown

import (
	"sync"
	"time"
)

// Cooldown tracks whether enough time has elapsed since the last activation.
type Cooldown struct {
	mu       sync.Mutex
	interval time.Duration
	lastFire time.Time
}

// New returns a Cooldown with the given minimum interval between activations.
// If interval is <= 0 it defaults to 1 second.
func New(interval time.Duration) *Cooldown {
	if interval <= 0 {
		interval = time.Second
	}
	return &Cooldown{interval: interval}
}

// Allow reports whether the cooldown period has elapsed since the last call
// to Allow that returned true. If it has, the internal timer is reset.
func (c *Cooldown) Allow() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	if c.lastFire.IsZero() || now.Sub(c.lastFire) >= c.interval {
		c.lastFire = now
		return true
	}
	return false
}

// Remaining returns the duration until the next activation is allowed.
// Returns zero if the cooldown has already elapsed.
func (c *Cooldown) Remaining() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lastFire.IsZero() {
		return 0
	}
	elapsed := time.Since(c.lastFire)
	if elapsed >= c.interval {
		return 0
	}
	return c.interval - elapsed
}

// Reset clears the cooldown state so the next call to Allow succeeds immediately.
func (c *Cooldown) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastFire = time.Time{}
}
