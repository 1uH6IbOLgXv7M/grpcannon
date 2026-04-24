// Package headroom estimates the remaining capacity of the load generator
// by comparing the current concurrency and error rate against configured
// maximums. A Headroom value of 1.0 means fully available; 0.0 means
// saturated.
package headroom

import "sync"

// Config holds the parameters used to compute headroom.
type Config struct {
	// MaxConcurrency is the upper bound of allowed in-flight requests.
	MaxConcurrency int
	// ErrorRateThreshold is the error rate [0,1] above which headroom is zero.
	ErrorRateThreshold float64
}

// Estimator computes a normalised headroom score in [0.0, 1.0].
type Estimator struct {
	mu  sync.Mutex
	cfg Config

	current    int
	errorRate  float64
}

// New returns an Estimator with the given Config.
// MaxConcurrency is clamped to at least 1.
// ErrorRateThreshold is clamped to [0, 1].
func New(cfg Config) *Estimator {
	if cfg.MaxConcurrency < 1 {
		cfg.MaxConcurrency = 1
	}
	if cfg.ErrorRateThreshold < 0 {
		cfg.ErrorRateThreshold = 0
	}
	if cfg.ErrorRateThreshold > 1 {
		cfg.ErrorRateThreshold = 1
	}
	return &Estimator{cfg: cfg}
}

// Update records the latest in-flight concurrency and error rate.
func (e *Estimator) Update(current int, errorRate float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.current = current
	e.errorRate = errorRate
}

// Score returns a value in [0.0, 1.0] representing available headroom.
// A score of 1.0 means no load; 0.0 means fully saturated or error-rate
// threshold exceeded.
func (e *Estimator) Score() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.cfg.ErrorRateThreshold > 0 && e.errorRate >= e.cfg.ErrorRateThreshold {
		return 0
	}

	ratio := float64(e.current) / float64(e.cfg.MaxConcurrency)
	if ratio > 1 {
		ratio = 1
	}
	return 1 - ratio
}

// Available reports whether there is any headroom remaining (Score > 0).
func (e *Estimator) Available() bool {
	return e.Score() > 0
}
