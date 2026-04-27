// Package leaky implements a leaky-bucket rate limiter that smooths
// bursts by draining at a fixed rate, dropping requests when the bucket
// is full.
package leaky

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrDropped is returned by Acquire when the bucket is full and the
// incoming request must be shed.
var ErrDropped = errors.New("leaky: bucket full, request dropped")

// Bucket is a leaky-bucket rate limiter.
type Bucket struct {
	mu       sync.Mutex
	current  float64
	capacity float64
	rate     float64 // tokens drained per second
	last     time.Time
}

// New returns a Bucket that drains at ratePerSec tokens per second and
// holds at most capacity tokens. Panics if either argument is <= 0.
func New(capacity, ratePerSec float64) *Bucket {
	if capacity <= 0 {
		panic("leaky: capacity must be > 0")
	}
	if ratePerSec <= 0 {
		panic("leaky: ratePerSec must be > 0")
	}
	return &Bucket{
		capacity: capacity,
		rate:     ratePerSec,
		last:     time.Now(),
	}
}

// Acquire attempts to add one token to the bucket. If the bucket is
// full ErrDropped is returned immediately. If ctx is already cancelled
// the context error is returned without modifying the bucket.
func (b *Bucket) Acquire(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.last).Seconds()
	b.last = now

	// Drain tokens that have leaked out since last call.
	b.current -= elapsed * b.rate
	if b.current < 0 {
		b.current = 0
	}

	if b.current >= b.capacity {
		return ErrDropped
	}
	b.current++
	return nil
}

// InFlight returns the approximate number of tokens currently held in
// the bucket.
func (b *Bucket) InFlight() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.current
}
