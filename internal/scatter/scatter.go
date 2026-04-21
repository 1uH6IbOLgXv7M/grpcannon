// Package scatter distributes snapshots across multiple named sinks,
// routing each snapshot to the sink whose tag key matches a label on the
// snapshot. Snapshots that match no registered sink are forwarded to an
// optional fallback sink.
package scatter

import "sync"

// Sink is any value that can receive a snapshot value.
type Sink[T any] interface {
	Send(T)
}

// Router holds a set of named sinks and an optional fallback.
type Router[T any] struct {
	mu       sync.RWMutex
	sinks    map[string]Sink[T]
	fallback Sink[T]
	key      func(T) string
}

// New creates a Router that extracts a routing key from each value using
// keyFn. If keyFn returns an empty string, or no sink is registered for
// the returned key, the value is forwarded to the fallback sink (if any).
func New[T any](keyFn func(T) string) *Router[T] {
	if keyFn == nil {
		panic("scatter: keyFn must not be nil")
	}
	return &Router[T]{
		sinks: make(map[string]Sink[T]),
		key:   keyFn,
	}
}

// Register adds or replaces the sink associated with name.
func (r *Router[T]) Register(name string, s Sink[T]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sinks[name] = s
}

// SetFallback sets the sink that receives values with no matching route.
func (r *Router[T]) SetFallback(s Sink[T]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fallback = s
}

// Send routes v to the registered sink for its key, or to the fallback.
func (r *Router[T]) Send(v T) {
	r.mu.RLock()
	s, ok := r.sinks[r.key(v)]
	fb := r.fallback
	r.mu.RUnlock()

	switch {
	case ok:
		s.Send(v)
	case fb != nil:
		fb.Send(v)
	}
}

// Len returns the number of registered (non-fallback) sinks.
func (r *Router[T]) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.sinks)
}
