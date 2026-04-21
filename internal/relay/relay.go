// Package relay provides a fan-out dispatcher that forwards snapshots to
// multiple downstream consumers, decoupling producers from consumers and
// absorbing back-pressure via per-sink bounded queues.
package relay

import (
	"context"
	"sync"

	"github.com/your-org/grpcannon/internal/metrics"
)

const defaultQueueDepth = 64

// Sink is any value that can receive a metrics snapshot.
type Sink interface {
	Receive(snap metrics.Snapshot)
}

// SinkFunc is a function adapter for Sink.
type SinkFunc func(metrics.Snapshot)

// Receive implements Sink.
func (f SinkFunc) Receive(snap metrics.Snapshot) { f(snap) }

// Relay fans out snapshots to a dynamic set of registered sinks.
// Each sink gets its own buffered channel so a slow consumer cannot
// block the producer or other consumers. Snapshots that cannot be
// enqueued immediately are dropped (non-blocking send).
type Relay struct {
	mu         sync.RWMutex
	sinks      map[string]sinkEntry
	queueDepth int
}

type sinkEntry struct {
	sink Sink
	ch   chan metrics.Snapshot
	cancel context.CancelFunc
}

// New returns a Relay with per-sink queue depth d.
// If d is less than 1 it defaults to 64.
func New(d int) *Relay {
	if d < 1 {
		d = defaultQueueDepth
	}
	return &Relay{
		sinks:      make(map[string]sinkEntry),
		queueDepth: d,
	}
}

// Register adds a named sink. If a sink with the same name already exists
// it is replaced; the old sink's goroutine is stopped gracefully.
func (r *Relay) Register(ctx context.Context, name string, s Sink) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if old, ok := r.sinks[name]; ok {
		old.cancel()
	}

	ch := make(chan metrics.Snapshot, r.queueDepth)
	ctx, cancel := context.WithCancel(ctx)

	r.sinks[name] = sinkEntry{sink: s, ch: ch, cancel: cancel}

	go func() {
		defer cancel()
		for {
			select {
			case snap, ok := <-ch:
				if !ok {
					return
				}
				s.Receive(snap)
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Unregister removes the named sink and stops its dispatch goroutine.
func (r *Relay) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if e, ok := r.sinks[name]; ok {
		e.cancel()
		delete(r.sinks, name)
	}
}

// Send dispatches snap to all registered sinks. Sends are non-blocking;
// if a sink's queue is full the snapshot is silently dropped for that sink.
func (r *Relay) Send(snap metrics.Snapshot) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, e := range r.sinks {
		select {
		case e.ch <- snap:
		default:
			// drop — sink is too slow
		}
	}
}

// Len returns the number of currently registered sinks.
func (r *Relay) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.sinks)
}
