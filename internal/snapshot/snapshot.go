// Package snapshot captures a point-in-time view of load test progress.
package snapshot

import (
	"time"

	"github.com/nicklaw5/grpcannon/internal/metrics"
)

// Snapshot holds a point-in-time summary of a running or completed load test.
type Snapshot struct {
	Timestamp  time.Time
	Elapsed    time.Duration
	Total      int64
	Errors     int64
	RPS        float64
	P50        time.Duration
	P95        time.Duration
	P99        time.Duration
	Mean       time.Duration
}

// Collector periodically samples metrics and stores snapshots.
type Collector struct {
	start    time.Time
	recorder *metrics.Recorder
	snaps    []Snapshot
}

// NewCollector creates a Collector backed by the given Recorder.
func NewCollector(r *metrics.Recorder) *Collector {
	return &Collector{
		start:    time.Now(),
		recorder: r,
	}
}

// Capture records a new snapshot from the current recorder state.
func (c *Collector) Capture() Snapshot {
	sum := c.recorder.Snapshot()
	elapsed := time.Since(c.start)

	var rps float64
	if secs := elapsed.Seconds(); secs > 0 {
		rps = float64(sum.Total) / secs
	}

	s := Snapshot{
		Timestamp: time.Now(),
		Elapsed:   elapsed,
		Total:     sum.Total,
		Errors:    sum.Errors,
		RPS:       rps,
		P50:       sum.P50,
		P95:       sum.P95,
		P99:       sum.P99,
		Mean:      sum.Mean,
	}
	c.snaps = append(c.snaps, s)
	return s
}

// All returns all captured snapshots in order.
func (c *Collector) All() []Snapshot {
	out := make([]Snapshot, len(c.snaps))
	copy(out, c.snaps)
	return out
}

// Latest returns the most recent snapshot, or a zero value if none exist.
func (c *Collector) Latest() Snapshot {
	if len(c.snaps) == 0 {
		return Snapshot{}
	}
	return c.snaps[len(c.snaps)-1]
}
