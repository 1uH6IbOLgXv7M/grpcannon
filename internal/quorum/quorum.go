// Package quorum decides whether a load test run has met a minimum
// success threshold before the caller proceeds or aborts.
package quorum

import (
	"errors"
	"fmt"
	"sync/atomic"
)

// ErrBelowThreshold is returned when the success ratio falls below the
// configured minimum.
var ErrBelowThreshold = errors.New("quorum: success ratio below threshold")

// Quorum tracks successes and failures and exposes a Check method that
// returns an error when the success ratio drops below MinRatio.
type Quorum struct {
	minRatio  float64
	minTotal  int64
	successes atomic.Int64
	failures  atomic.Int64
}

// New returns a Quorum that requires at least minTotal observations before
// enforcing minRatio (0–1).
func New(minRatio float64, minTotal int64) *Quorum {
	if minRatio < 0 {
		minRatio = 0
	}
	if minRatio > 1 {
		minRatio = 1
	}
	if minTotal < 1 {
		minTotal = 1
	}
	return &Quorum{minRatio: minRatio, minTotal: minTotal}
}

// RecordSuccess increments the success counter.
func (q *Quorum) RecordSuccess() { q.successes.Add(1) }

// RecordFailure increments the failure counter.
func (q *Quorum) RecordFailure() { q.failures.Add(1) }

// Ratio returns the current success ratio. Returns 1.0 when no observations
// have been recorded yet.
func (q *Quorum) Ratio() float64 {
	s := q.successes.Load()
	f := q.failures.Load()
	total := s + f
	if total == 0 {
		return 1.0
	}
	return float64(s) / float64(total)
}

// Total returns the total number of observations recorded.
func (q *Quorum) Total() int64 {
	return q.successes.Load() + q.failures.Load()
}

// Check returns ErrBelowThreshold when enough observations have been
// collected and the success ratio is below the configured minimum.
func (q *Quorum) Check() error {
	if q.Total() < q.minTotal {
		return nil
	}
	ratio := q.Ratio()
	if ratio < q.minRatio {
		return fmt.Errorf("%w: got %.4f want >= %.4f", ErrBelowThreshold, ratio, q.minRatio)
	}
	return nil
}
