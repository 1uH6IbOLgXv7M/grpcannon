// Package budget provides an error-budget tracker that opens a gate
// once the ratio of failed requests exceeds a configured threshold.
package budget

import (
	"errors"
	"sync"
)

// ErrExceeded is returned by Acquire when the error budget is spent.
var ErrExceeded = errors.New("budget: error budget exceeded")

// Budget tracks successes and failures and blocks new requests once
// the failure ratio exceeds the configured threshold.
type Budget struct {
	mu        sync.Mutex
	threshold float64 // 0‥1
	total     int64
	failures  int64
}

// New returns a Budget that trips when failure/total >= threshold.
// threshold must be in (0, 1]; values outside that range are clamped.
func New(threshold float64) *Budget {
	if threshold <= 0 {
		threshold = 0.01
	}
	if threshold > 1 {
		threshold = 1
	}
	return &Budget{threshold: threshold}
}

// Record registers the outcome of a single request.
func (b *Budget) Record(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.total++
	if err != nil {
		b.failures++
	}
}

// Allow returns ErrExceeded when the current failure ratio is at or
// above the threshold and at least one request has been recorded.
func (b *Budget) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.total == 0 {
		return nil
	}
	if float64(b.failures)/float64(b.total) >= b.threshold {
		return ErrExceeded
	}
	return nil
}

// Ratio returns the current failure ratio.
func (b *Budget) Ratio() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.total == 0 {
		return 0
	}
	return float64(b.failures) / float64(b.total)
}

// Reset clears all counters.
func (b *Budget) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.total = 0
	b.failures = 0
}
