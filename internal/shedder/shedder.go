// Package shedder implements load-based request shedding using a sliding
// window error rate combined with an in-flight request ceiling.
package shedder

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/nickbadlose/grpcannon/internal/window"
)

// ErrShed is returned when a request is shed due to load.
var ErrShed = errors.New("shedder: request shed")

// Config holds tuning parameters for the Shedder.
type Config struct {
	// MaxInFlight is the hard ceiling on concurrent requests. Zero means unlimited.
	MaxInFlight int64
	// ErrorRateThreshold in [0,1] above which new requests are shed. Zero disables.
	ErrorRateThreshold float64
	// WindowSize is the rolling window used to compute the error rate.
	WindowSize time.Duration
}

// Shedder decides whether to accept or shed an incoming request.
type Shedder struct {
	cfg      Config
	inFlight atomic.Int64
	win       *window.Window
}

// New returns a ready-to-use Shedder.
func New(cfg Config) *Shedder {
	if cfg.WindowSize <= 0 {
		cfg.WindowSize = 5 * time.Second
	}
	return &Shedder{
		cfg: cfg,
		win:  window.New(cfg.WindowSize),
	}
}

// Acquire attempts to reserve a slot for an incoming request.
// Returns ErrShed if the request should be dropped.
func (s *Shedder) Acquire() error {
	if s.cfg.MaxInFlight > 0 && s.inFlight.Load() >= s.cfg.MaxInFlight {
		s.win.Add(0, 1)
		return ErrShed
	}
	if s.cfg.ErrorRateThreshold > 0 && s.win.ErrorRate() >= s.cfg.ErrorRateThreshold {
		s.win.Add(0, 1)
		return ErrShed
	}
	s.inFlight.Add(1)
	return nil
}

// Release decrements the in-flight counter and records the outcome.
func (s *Shedder) Release(err error) {
	s.inFlight.Add(-1)
	if err != nil {
		s.win.Add(1, 1)
	} else {
		s.win.Add(0, 1)
	}
}

// InFlight returns the current number of in-flight requests.
func (s *Shedder) InFlight() int64 { return s.inFlight.Load() }

// ErrorRate returns the current rolling error rate.
func (s *Shedder) ErrorRate() float64 { return s.win.ErrorRate() }
