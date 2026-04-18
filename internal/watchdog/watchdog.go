// Package watchdog monitors error rates and cancels the load test when a
// configurable threshold is breached, preventing runaway failures.
package watchdog

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrThresholdExceeded is returned when the error rate exceeds the limit.
var ErrThresholdExceeded = errors.New("watchdog: error rate threshold exceeded")

// Config holds watchdog parameters.
type Config struct {
	// Threshold is the maximum tolerated error rate in [0,1].
	Threshold float64
	// Window is how far back to look when computing the error rate.
	Window time.Duration
	// MinRequests is the minimum number of requests before enforcement begins.
	MinRequests int64
}

// Watchdog polls an error-rate source and cancels a context when the threshold
// is breached.
type Watchdog struct {
	cfg    Config
	source Source
	once   sync.Once
}

// Source provides live error-rate statistics.
type Source interface {
	ErrorRate() float64
	Total() int64
}

// New returns a Watchdog configured with cfg.
func New(cfg Config, src Source) *Watchdog {
	if cfg.Window <= 0 {
		cfg.Window = 5 * time.Second
	}
	if cfg.MinRequests <= 0 {
		cfg.MinRequests = 10
	}
	return &Watchdog{cfg: cfg, source: src}
}

// Run starts the watchdog loop. It cancels cancel when the threshold is
// breached and returns ErrThresholdExceeded. It returns nil when ctx is done.
func (w *Watchdog) Run(ctx context.Context, cancel context.CancelCauseFunc) error {
	ticker := time.NewTicker(w.cfg.Window / 5)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if w.source.Total() < w.cfg.MinRequests {
				continue
			}
			if w.source.ErrorRate() > w.cfg.Threshold {
				w.once.Do(func() { cancel(ErrThresholdExceeded) })
				return ErrThresholdExceeded
			}
		}
	}
}
