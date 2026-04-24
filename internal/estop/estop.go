// Package estop provides an emergency-stop mechanism that trips once a
// configurable error-rate threshold is exceeded within a sliding window,
// preventing further requests from being dispatched until the operator
// explicitly resets the stop.
package estop

import (
	"errors"
	"sync"
	"sync/atomic"
)

// ErrTripped is returned by Allow when the emergency stop has been tripped.
var ErrTripped = errors.New("estop: emergency stop tripped")

// EStop trips once the error rate inside a sliding window exceeds Threshold.
type EStop struct {
	mu        sync.Mutex
	threshold float64 // 0–1
	window    int64   // total observations in window
	errors    int64   // error observations in window
	tripped   atomic.Bool
}

// New returns an EStop that trips when the error rate exceeds threshold (0–1).
// A threshold of 0 means any single error trips the stop.
func New(threshold float64) *EStop {
	if threshold < 0 {
		threshold = 0
	}
	if threshold > 1 {
		threshold = 1
	}
	return &EStop{threshold: threshold}
}

// Allow returns ErrTripped if the stop has been tripped, otherwise nil.
func (e *EStop) Allow() error {
	if e.tripped.Load() {
		return ErrTripped
	}
	return nil
}

// RecordSuccess records a successful observation and evaluates the threshold.
func (e *EStop) RecordSuccess() {
	e.mu.Lock()
	e.window++
	e.mu.Unlock()
}

// RecordFailure records a failed observation and trips the stop if the
// error rate now exceeds the configured threshold.
func (e *EStop) RecordFailure() {
	e.mu.Lock()
	e.window++
	e.errors++
	rate := float64(e.errors) / float64(e.window)
	should := rate > e.threshold
	e.mu.Unlock()
	if should {
		e.tripped.Store(true)
	}
}

// Tripped reports whether the stop has been tripped.
func (e *EStop) Tripped() bool { return e.tripped.Load() }

// Reset clears the tripped state and resets the internal counters.
func (e *EStop) Reset() {
	e.mu.Lock()
	e.window = 0
	e.errors = 0
	e.mu.Unlock()
	e.tripped.Store(false)
}

// ErrorRate returns the current error rate observed since the last Reset.
func (e *EStop) ErrorRate() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.window == 0 {
		return 0
	}
	return float64(e.errors) / float64(e.window)
}
