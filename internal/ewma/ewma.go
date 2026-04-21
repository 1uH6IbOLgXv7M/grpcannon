// Package ewma provides an exponentially weighted moving average.
// It is useful for tracking smoothed rate or latency signals over time
// without retaining a full history of observations.
package ewma

import (
	"math"
	"sync"
)

// DefaultDecay is the smoothing factor used when none is specified.
// A value of 0.1 weights recent observations more heavily.
const DefaultDecay = 0.1

// EWMA is a thread-safe exponentially weighted moving average.
type EWMA struct {
	mu    sync.Mutex
	alpha float64
	value float64
	init  bool
}

// New returns an EWMA with the given decay factor alpha in (0, 1].
// Values outside that range are clamped to DefaultDecay.
func New(alpha float64) *EWMA {
	if alpha <= 0 || alpha > 1 {
		alpha = DefaultDecay
	}
	return &EWMA{alpha: alpha}
}

// Add incorporates a new observation into the moving average.
func (e *EWMA) Add(v float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.init {
		e.value = v
		e.init = true
		return
	}
	e.value = e.alpha*v + (1-e.alpha)*e.value
}

// Value returns the current smoothed average.
// Returns 0 if no observations have been added.
func (e *EWMA) Value() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.value
}

// Reset clears the EWMA back to its initial state.
func (e *EWMA) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.value = 0
	e.init = false
}

// Rate returns the EWMA value rounded to the given number of decimal places.
func (e *EWMA) Rate(decimals int) float64 {
	v := e.Value()
	if decimals < 0 {
		decimals = 0
	}
	shift := math.Pow(10, float64(decimals))
	return math.Round(v*shift) / shift
}
