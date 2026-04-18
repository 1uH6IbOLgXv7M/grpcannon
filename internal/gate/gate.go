// Package gate provides a concurrency gate that limits the number of
// concurrent operations allowed to proceed, blocking callers until a slot
// is available or the context is cancelled.
package gate

import (
	"context"
	"errors"
	"sync"
)

// ErrClosed is returned when Wait is called on a closed Gate.
var ErrClosed = errors.New("gate: closed")

// Gate controls how many goroutines may proceed concurrently.
type Gate struct {
	mu     sync.Mutex
	sem    chan struct{}
	closed bool
}

// New creates a Gate that allows at most n concurrent operations.
// If n <= 0, the gate is unbounded.
func New(n int) *Gate {
	if n <= 0 {
		return &Gate{}
	}
	g := &Gate{sem: make(chan struct{}, n)}
	for i := 0; i < n; i++ {
		g.sem <- struct{}{}
	}
	return g
}

// Wait blocks until a slot is available, the context is done, or the gate
// is closed. Callers must call Done after their operation completes.
func (g *Gate) Wait(ctx context.Context) error {
	g.mu.Lock()
	if g.closed {
		g.mu.Unlock()
		return ErrClosed
	}
	g.mu.Unlock()

	if g.sem == nil {
		return nil
	}
	select {
	case <-g.sem:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Done releases a slot back to the gate.
func (g *Gate) Done() {
	if g.sem == nil {
		return
	}
	g.sem <- struct{}{}
}

// Close permanently closes the gate. Subsequent calls to Wait return ErrClosed.
func (g *Gate) Close() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.closed = true
}
