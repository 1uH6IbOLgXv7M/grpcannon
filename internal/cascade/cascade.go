// Package cascade provides a sequential failure detector that trips after
// a configurable number of consecutive errors and resets on success.
package cascade

import "sync"

// ErrOpen is returned when the cascade detector is tripped.
type ErrOpen struct{}

func (e ErrOpen) Error() string { return "cascade: consecutive error threshold exceeded" }

// Detector trips after Threshold consecutive failures and resets on any success.
type Detector struct {
	mu        sync.Mutex
	threshold int
	consec    int
	tripped   bool
}

// New returns a Detector that trips after threshold consecutive failures.
// threshold must be >= 1; values below 1 are clamped to 1.
func New(threshold int) *Detector {
	if threshold < 1 {
		threshold = 1
	}
	return &Detector{threshold: threshold}
}

// RecordSuccess resets the consecutive failure counter and clears the tripped state.
func (d *Detector) RecordSuccess() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.consec = 0
	d.tripped = false
}

// RecordFailure increments the consecutive failure counter.
// Once the threshold is reached the detector is tripped.
func (d *Detector) RecordFailure() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.consec++
	if d.consec >= d.threshold {
		d.tripped = true
	}
}

// Allow returns nil when the detector is healthy, or ErrOpen when tripped.
func (d *Detector) Allow() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.tripped {
		return ErrOpen{}
	}
	return nil
}

// Consecutive returns the current consecutive failure count.
func (d *Detector) Consecutive() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.consec
}

// Tripped reports whether the detector is currently tripped.
func (d *Detector) Tripped() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.tripped
}
