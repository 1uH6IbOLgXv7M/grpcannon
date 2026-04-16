// Package drain provides graceful shutdown helpers that wait for
// in-flight requests to complete before the process exits.
package drain

import (
	"context"
	"sync"
	"time"
)

// Drainer tracks active requests and blocks until all finish or the
// deadline is exceeded.
type Drainer struct {
	mu      sync.Mutex
	wg      sync.WaitGroup
	closed  bool
}

// New returns a ready-to-use Drainer.
func New() *Drainer {
	return &Drainer{}
}

// Acquire marks one request as in-flight. It returns false if the
// Drainer has already been closed, meaning no new work should start.
func (d *Drainer) Acquire() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return false
	}
	d.wg.Add(1)
	return true
}

// Release marks one in-flight request as complete.
func (d *Drainer) Release() {
	d.wg.Done()
}

// Drain closes the Drainer to new acquisitions and blocks until all
// in-flight requests finish or ctx is cancelled. It returns
// context.DeadlineExceeded / context.Canceled on timeout.
func (d *Drainer) Drain(ctx context.Context) error {
	d.mu.Lock()
	d.closed = true
	d.mu.Unlock()

	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// DrainTimeout is a convenience wrapper around Drain with a fixed timeout.
func (d *Drainer) DrainTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return d.Drain(ctx)
}
