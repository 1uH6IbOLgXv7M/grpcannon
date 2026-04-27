// Package tokenring implements a round-robin token distributor that assigns
// request slots to workers in a fair, rotating order. It is useful when a
// fixed number of upstream endpoints must be exercised evenly during a load
// test run.
package tokenring

import (
	"errors"
	"sync"
	"sync/atomic"
)

// ErrEmpty is returned when the ring contains no slots.
var ErrEmpty = errors.New("tokenring: ring is empty")

// Ring distributes integer tokens in round-robin order across a fixed set of
// slots. It is safe for concurrent use.
type Ring struct {
	mu    sync.RWMutex
	slots []int
	cursor atomic.Uint64
}

// New creates a Ring pre-populated with tokens [0, size).
// size must be >= 1; if it is < 1 it is clamped to 1.
func New(size int) *Ring {
	if size < 1 {
		size = 1
	}
	slots := make([]int, size)
	for i := range slots {
		slots[i] = i
	}
	return &Ring{slots: slots}
}

// Next returns the next token in round-robin order.
// It returns ErrEmpty only when the ring has been drained via Reset(0).
func (r *Ring) Next() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.slots) == 0 {
		return 0, ErrEmpty
	}
	idx := r.cursor.Add(1) - 1
	return r.slots[int(idx)%len(r.slots)], nil
}

// Len returns the current number of slots in the ring.
func (r *Ring) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.slots)
}

// Reset replaces the ring contents with tokens [0, size).
// A size < 1 clears the ring (Next will return ErrEmpty).
func (r *Ring) Reset(size int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if size < 1 {
		r.slots = nil
		r.cursor.Store(0)
		return
	}
	slots := make([]int, size)
	for i := range slots {
		slots[i] = i
	}
	r.slots = slots
	r.cursor.Store(0)
}
