// Package slope computes the rate of change (slope) of a metric over a
// sliding window of observations. It uses a simple linear regression over
// the most recent samples to estimate whether a value is rising, falling,
// or stable. This is useful for adaptive controllers that need to react to
// trends rather than instantaneous values.
package slope

import "sync"

// Estimator tracks a fixed-size window of (x, y) samples and exposes the
// least-squares slope of the series. x is the sample index (monotonically
// increasing); y is the observed metric value.
type Estimator struct {
	mu      sync.Mutex
	window  int
	xs      []float64
	ys      []float64
	cursor  int
	full    bool
}

// New returns an Estimator that retains the last window samples.
// If window is less than 2 it is clamped to 2.
func New(window int) *Estimator {
	if window < 2 {
		window = 2
	}
	return &Estimator{
		window: window,
		xs:     make([]float64, window),
		ys:     make([]float64, window),
	}
}

// Add records a new observation. x should be a monotonically increasing
// value (e.g. elapsed seconds or a sequence number); y is the metric.
func (e *Estimator) Add(x, y float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.xs[e.cursor] = x
	e.ys[e.cursor] = y
	e.cursor = (e.cursor + 1) % e.window
	if e.cursor == 0 {
		e.full = true
	}
}

// Slope returns the least-squares slope over the current window.
// Returns 0 if fewer than 2 samples have been recorded.
func (e *Estimator) Slope() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()

	n := e.window
	if !e.full {
		n = e.cursor
	}
	if n < 2 {
		return 0
	}

	var sumX, sumY, sumXY, sumXX float64
	for i := 0; i < n; i++ {
		x := e.xs[i]
		y := e.ys[i]
		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}
	fn := float64(n)
	denom := fn*sumXX - sumX*sumX
	if denom == 0 {
		return 0
	}
	return (fn*sumXY - sumX*sumY) / denom
}

// Count returns the number of observations recorded so far, capped at window.
func (e *Estimator) Count() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.full {
		return e.window
	}
	return e.cursor
}
