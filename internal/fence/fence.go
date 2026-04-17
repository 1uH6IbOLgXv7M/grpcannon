// Package fence provides a reusable one-shot barrier that allows multiple
// goroutines to wait until a condition is signalled, then proceeds immediately
// for all subsequent callers.
package fence

import "sync"

// Fence is a one-shot barrier. Once opened it stays open forever.
type Fence struct {
	mu     sync.Mutex
	open   bool
	notify chan struct{}
}

// New returns a closed Fence.
func New() *Fence {
	return &Fence{notify: make(chan struct{})}
}

// Open unblocks all current and future calls to Wait.
func (f *Fence) Open() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.open {
		f.open = true
		close(f.notify)
	}
}

// IsOpen reports whether the fence has been opened.
func (f *Fence) IsOpen() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.open
}

// Wait blocks until Open is called or ctx is cancelled.
// Returns nil when the fence opens, ctx.Err() otherwise.
func (f *Fence) Wait(ctx context.Context) error {
	select {
	case <-f.notify:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
