// Package metrics provides latency histogram and result aggregation
// for gRPC load test runs.
package metrics

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

// Recorder collects latency samples and error counts in a thread-safe manner.
type Recorder struct {
	mu       sync.Mutex
	latencies []time.Duration
	errorCount int
	total      int
}

// NewRecorder returns an initialised Recorder.
func NewRecorder() *Recorder {
	return &Recorder{}
}

// Record adds a single observation. If err is non-nil the sample is counted
// as an error and the latency is still recorded.
func (r *Recorder) Record(d time.Duration, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.latencies = append(r.latencies, d)
	r.total++
	if err != nil {
		r.errorCount++
	}
}

// Snapshot returns an immutable Summary computed from all recorded samples.
func (r *Recorder) Snapshot() Summary {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.latencies) == 0 {
		return Summary{Total: r.total, Errors: r.errorCount}
	}

	sorted := make([]time.Duration, len(r.latencies))
	copy(sorted, r.latencies)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	var sum time.Duration
	for _, v := range sorted {
		sum += v
	}

	return Summary{
		Total:   r.total,
		Errors:  r.errorCount,
		Min:     sorted[0],
		Max:     sorted[len(sorted)-1],
		Mean:    time.Duration(int64(sum) / int64(len(sorted))),
		P50:     percentile(sorted, 50),
		P95:     percentile(sorted, 95),
		P99:     percentile(sorted, 99),
	}
}

// Summary holds aggregated statistics for a completed run.
type Summary struct {
	Total  int
	Errors int
	Min    time.Duration
	Max    time.Duration
	Mean   time.Duration
	P50    time.Duration
	P95    time.Duration
	P99    time.Duration
}

// String returns a human-readable representation of the Summary.
func (s Summary) String() string {
	successRate := 0.0
	if s.Total > 0 {
		successRate = float64(s.Total-s.Errors) / float64(s.Total) * 100
	}
	return fmt.Sprintf(
		"Total: %d | Errors: %d (%.1f%% success)\n"+
			"Min: %v | Mean: %v | Max: %v\n"+
			"P50: %v | P95: %v | P99: %v",
		s.Total, s.Errors, successRate,
		s.Min, s.Mean, s.Max,
		s.P50, s.P95, s.P99,
	)
}

// percentile returns the p-th percentile value from a pre-sorted slice.
func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(math.Ceil(p/100.0*float64(len(sorted)))) - 1
	if idx < 0 {
		idx = 0
	}
	return sorted[idx]
}
