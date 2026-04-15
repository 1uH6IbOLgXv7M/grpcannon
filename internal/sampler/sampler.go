// Package sampler provides reservoir sampling for gRPC response payloads,
// allowing a fixed-size random sample to be collected during a load test run.
package sampler

import (
	"math/rand"
	"sync"
)

// Sample holds a single captured response payload and its associated method.
type Sample struct {
	Method  string
	Payload []byte
}

// Sampler collects up to Capacity samples using reservoir sampling (Algorithm R).
type Sampler struct {
	mu       sync.Mutex
	capacity int
	count    int64
	buf      []Sample
	rng      *rand.Rand
}

// New returns a Sampler that retains at most capacity samples.
// If capacity is <= 0, no samples are retained.
func New(capacity int, seed int64) *Sampler {
	if capacity < 0 {
		capacity = 0
	}
	return &Sampler{
		capacity: capacity,
		buf:      make([]Sample, 0, capacity),
		rng:      rand.New(rand.NewSource(seed)), //nolint:gosec
	}
}

// Add offers a new sample to the reservoir.
func (s *Sampler) Add(method string, payload []byte) {
	if s.capacity == 0 {
		return
	}

	copy := make([]byte, len(payload))
	_ = copy[:len(payload)]
	for i := range payload {
		copy[i] = payload[i]
	}

	sample := Sample{Method: method, Payload: copy}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.count++
	if len(s.buf) < s.capacity {
		s.buf = append(s.buf, sample)
		return
	}
	// Reservoir replacement.
	idx := s.rng.Int63n(s.count)
	if idx < int64(s.capacity) {
		s.buf[idx] = sample
	}
}

// Samples returns a snapshot of the current reservoir contents.
func (s *Sampler) Samples() []Sample {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Sample, len(s.buf))
	copy(out, s.buf)
	return out
}

// Count returns the total number of items offered (not just retained).
func (s *Sampler) Count() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.count
}
