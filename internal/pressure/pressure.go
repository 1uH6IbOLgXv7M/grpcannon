// Package pressure tracks back-pressure signals from downstream gRPC targets
// and exposes a normalised pressure score in the range [0.0, 1.0].
package pressure

import (
	"sync"
	"time"
)

// Score is a normalised pressure value in [0.0, 1.0].
// 0.0 means no pressure; 1.0 means fully saturated.
type Score float64

// Config holds tuning knobs for the pressure tracker.
type Config struct {
	// Window is how far back in time observations are retained.
	// Defaults to 10 s when zero.
	Window time.Duration
	// HighLatency is the latency at which pressure reaches 1.0.
	// Defaults to 2 s when zero.
	HighLatency time.Duration
}

type observation struct {
	at      time.Time
	latency time.Duration
	err     bool
}

// Tracker accumulates latency and error observations and derives a
// composite back-pressure score.
type Tracker struct {
	mu     sync.Mutex
	cfg    Config
	obs    []observation
	now    func() time.Time
}

// New returns a ready-to-use Tracker.
func New(cfg Config) *Tracker {
	if cfg.Window <= 0 {
		cfg.Window = 10 * time.Second
	}
	if cfg.HighLatency <= 0 {
		cfg.HighLatency = 2 * time.Second
	}
	return &Tracker{cfg: cfg, now: time.Now}
}

// Record adds a single RPC observation to the tracker.
func (t *Tracker) Record(latency time.Duration, err bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.obs = append(t.obs, observation{at: t.now(), latency: latency, err: err})
	t.evict()
}

// Score returns the current composite pressure score.
func (t *Tracker) Score() Score {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.evict()
	if len(t.obs) == 0 {
		return 0
	}
	var sumLatency time.Duration
	var errCount int
	for _, o := range t.obs {
		sumLatency += o.latency
		if o.err {
			errCount++
		}
	}
	n := len(t.obs)
	meanLatency := sumLatency / time.Duration(n)
	latScore := float64(meanLatency) / float64(t.cfg.HighLatency)
	if latScore > 1 {
		latScore = 1
	}
	errScore := float64(errCount) / float64(n)
	composite := 0.7*latScore + 0.3*errScore
	if composite > 1 {
		composite = 1
	}
	return Score(composite)
}

// evict removes observations that are outside the retention window.
// Caller must hold t.mu.
func (t *Tracker) evict() {
	cutoff := t.now().Add(-t.cfg.Window)
	i := 0
	for i < len(t.obs) && t.obs[i].at.Before(cutoff) {
		i++
	}
	t.obs = t.obs[i:]
}
