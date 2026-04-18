// Package tee fans out metrics snapshots to multiple consumers.
package tee

import (
	"context"
	"sync"

	"github.com/example/grpcannon/internal/metrics"
)

// Sink is any value that can receive a metrics snapshot.
type Sink interface {
	Write(snap metrics.Snapshot)
}

// Tee fans a stream of snapshots out to a set of registered sinks.
type Tee struct {
	mu    sync.RWMutex
	sinks []Sink
}

// New returns an empty Tee.
func New(sinks ...Sink) *Tee {
	return &Tee{sinks: append([]Sink{}, sinks...)}
}

// Add registers an additional sink.
func (t *Tee) Add(s Sink) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.sinks = append(t.sinks, s)
}

// Send delivers snap to every registered sink concurrently and waits for
// all deliveries to complete before returning.
func (t *Tee) Send(snap metrics.Snapshot) {
	t.mu.RLock()
	copy := append([]Sink{}, t.sinks...)
	t.mu.RUnlock()

	var wg sync.WaitGroup
	for _, s := range copy {
		wg.Add(1)
		go func(s Sink) {
			defer wg.Done()
			s.Write(snap)
		}(s)
	}
	wg.Wait()
}

// Run reads snapshots from ch and fans them out until ctx is done or ch is
// closed.
func (t *Tee) Run(ctx context.Context, ch <-chan metrics.Snapshot) {
	for {
		select {
		case <-ctx.Done():
			return
		case snap, ok := <-ch:
			if !ok {
				return
			}
			t.Send(snap)
		}
	}
}
