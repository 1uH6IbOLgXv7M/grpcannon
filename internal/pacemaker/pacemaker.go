// Package pacemaker dynamically adjusts the request rate based on observed
// latency, keeping the p99 latency below a configurable target.
package pacemaker

import (
	"sync"
	"time"
)

// Config holds tuning parameters for the pacemaker.
type Config struct {
	// TargetP99 is the desired p99 latency ceiling.
	TargetP99 time.Duration
	// MinRPS is the floor for the computed rate (0 = unlimited floor).
	MinRPS float64
	// MaxRPS is the ceiling for the computed rate (0 = no ceiling).
	MaxRPS float64
	// StepFactor controls how aggressively the rate is adjusted (default 0.1).
	StepFactor float64
}

// Pacemaker watches p99 latency samples and emits a suggested RPS.
type Pacemaker struct {
	cfg    Config
	mu     sync.Mutex
	current float64
}

// New returns a Pacemaker initialised at maxRPS (or 100 when maxRPS is 0).
func New(cfg Config) *Pacemaker {
	if cfg.StepFactor <= 0 {
		cfg.StepFactor = 0.10
	}
	initial := cfg.MaxRPS
	if initial <= 0 {
		initial = 100
	}
	if cfg.MinRPS > 0 && initial < cfg.MinRPS {
		initial = cfg.MinRPS
	}
	return &Pacemaker{cfg: cfg, current: initial}
}

// Adjust recalculates the suggested RPS given the latest p99 observation and
// returns the new value.
func (p *Pacemaker) Adjust(p99 time.Duration) float64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cfg.TargetP99 <= 0 {
		return p.current
	}

	ratio := float64(p99) / float64(p.cfg.TargetP99)
	switch {
	case ratio > 1.0:
		// latency too high – slow down
		p.current *= (1 - p.cfg.StepFactor*ratio)
	case ratio < 0.9:
		// latency comfortably below target – speed up
		p.current *= (1 + p.cfg.StepFactor*(1-ratio))
	}

	if p.cfg.MinRPS > 0 && p.current < p.cfg.MinRPS {
		p.current = p.cfg.MinRPS
	}
	if p.cfg.MaxRPS > 0 && p.current > p.cfg.MaxRPS {
		p.current = p.cfg.MaxRPS
	}
	return p.current
}

// Current returns the most recently computed RPS without modifying it.
func (p *Pacemaker) Current() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.current
}
