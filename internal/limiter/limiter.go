// Package limiter provides a concurrency limiter that caps the number of
// goroutines executing simultaneously.
package limiter

import (
	"context"
	"errors"
)

// ErrLimitExceeded is returned when the limiter is full and the context is
// already done.
var ErrLimitExceeded = errors.New("limiter: concurrency limit exceeded")

// Limiter gates concurrent execution using a semaphore channel.
type Limiter struct {
	sem chan struct{}
}

// New returns a Limiter that allows at most n concurrent acquisitions.
// If n <= 0 the limiter is unbounded.
func New(n int) *Limiter {
	if n <= 0 {
		return &Limiter{}
	}
	return &Limiter{sem: make(chan struct{}, n)}
}

// Acquire blocks until a slot is available or ctx is done.
// Returns ErrLimitExceeded if the context expires while waiting.
func (l *Limiter) Acquire(ctx context.Context) error {
	if l.sem == nil {
		return nil
	}
	select {
	case l.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ErrLimitExceeded
	}
}

// Release frees a previously acquired slot.
func (l *Limiter) Release() {
	if l.sem == nil {
		return
	}
	select {
	case <-l.sem:
	default:
	}
}

// Available returns the number of free slots remaining.
// Returns -1 for an unbounded limiter.
func (l *Limiter) Available() int {
	if l.sem == nil {
		return -1
	}
	return cap(l.sem) - len(l.sem)
}
