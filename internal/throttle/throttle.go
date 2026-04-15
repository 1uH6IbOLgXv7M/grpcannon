// Package throttle provides a token-bucket style request throttler that
// integrates with the concurrency profile to dynamically adjust throughput
// as worker counts change during a load test.
package throttle

import (
	"context"
	"sync"
	"time"
)

// Throttle controls the rate at which requests are dispatched.
type Throttle struct {
	mu       sync.Mutex
	ticker   *time.Ticker
	tokens   chan struct{}
	stop     chan struct{}
	stopped  bool
}

// New creates a Throttle that emits at most rps tokens per second.
// If rps <= 0 the throttle is unlimited and Wait returns immediately.
func New(rps int) *Throttle {
	t := &Throttle{
		stop: make(chan struct{}),
	}
	if rps <= 0 {
		return t
	}
	interval := time.Second / time.Duration(rps)
	t.tokens = make(chan struct{}, rps)
	t.ticker = time.NewTicker(interval)
	go t.fill()
	return t
}

func (t *Throttle) fill() {
	for {
		select {
		case <-t.ticker.C:
			select {
			case t.tokens <- struct{}{}:
			default:
			}
		case <-t.stop:
			return
		}
	}
}

// Wait blocks until a token is available or ctx is cancelled.
// Returns ctx.Err() if the context is cancelled before a token arrives.
func (t *Throttle) Wait(ctx context.Context) error {
	if t.tokens == nil {
		return nil
	}
	select {
	case <-t.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop releases resources held by the Throttle.
// Calling Stop more than once is safe.
func (t *Throttle) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.stopped || t.ticker == nil {
		return
	}
	t.stopped = true
	t.ticker.Stop()
	close(t.stop)
}
