// Package overflow provides a fixed-capacity drop policy for excess requests.
// When the queue is full, new requests are dropped and counted.
package overflow

import "sync/atomic"

// Policy controls what happens when capacity is exceeded.
type Policy int

const (
	// Drop silently discards excess requests.
	Drop Policy = iota
)

// Queue is a bounded channel-backed queue that counts dropped items.
type Queue struct {
	ch      chan struct{}
	dropped atomic.Int64
	total   atomic.Int64
}

// New returns a Queue with the given capacity. Capacity must be >= 1.
func New(capacity int) *Queue {
	if capacity < 1 {
		capacity = 1
	}
	return &Queue{ch: make(chan struct{}, capacity)}
}

// Acquire attempts to claim a slot. Returns true if acquired, false if dropped.
func (q *Queue) Acquire() bool {
	q.total.Add(1)
	select {
	case q.ch <- struct{}{}:
		return true
	default:
		q.dropped.Add(1)
		return false
	}
}

// Release frees a previously acquired slot.
func (q *Queue) Release() {
	select {
	case <-q.ch:
	default:
	}
}

// Dropped returns the total number of dropped requests.
func (q *Queue) Dropped() int64 { return q.dropped.Load() }

// Total returns the total number of Acquire attempts.
func (q *Queue) Total() int64 { return q.total.Load() }

// InFlight returns the current number of acquired slots.
func (q *Queue) InFlight() int { return len(q.ch) }
