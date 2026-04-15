// Package histogram provides a fixed-bucket latency histogram for
// summarising request durations during a load test run.
package histogram

import (
	"fmt"
	"strings"
	"time"
)

// defaultBuckets are upper-bound boundaries in milliseconds.
var defaultBuckets = []float64{1, 2.5, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000}

// Histogram counts observations into fixed-width buckets.
type Histogram struct {
	buckets []float64 // upper bounds in ms
	counts  []int64
	overflow int64 // observations above the largest bucket
}

// New returns a Histogram using the default latency buckets.
func New() *Histogram {
	return NewWithBuckets(defaultBuckets)
}

// NewWithBuckets returns a Histogram with the supplied upper-bound values (ms).
func NewWithBuckets(buckets []float64) *Histogram {
	b := make([]float64, len(buckets))
	copy(b, buckets)
	return &Histogram{
		buckets: b,
		counts:  make([]int64, len(b)),
	}
}

// Observe records a single duration observation.
func (h *Histogram) Observe(d time.Duration) {
	ms := float64(d) / float64(time.Millisecond)
	for i, upper := range h.buckets {
		if ms <= upper {
			h.counts[i]++
			return
		}
	}
	h.overflow++
}

// Total returns the total number of observations recorded.
func (h *Histogram) Total() int64 {
	var n int64
	for _, c := range h.counts {
		n += c
	}
	return n + h.overflow
}

// String renders the histogram as a simple ASCII bar chart.
func (h *Histogram) String() string {
	total := h.Total()
	if total == 0 {
		return "(no observations)"
	}
	var sb strings.Builder
	for i, upper := range h.buckets {
		pct := float64(h.counts[i]) / float64(total) * 100
		bar := strings.Repeat("█", int(pct/5))
		fmt.Fprintf(&sb, "<= %6.1f ms | %-20s %6.1f%% (%d)\n", upper, bar, pct, h.counts[i])
	}
	if h.overflow > 0 {
		pct := float64(h.overflow) / float64(total) * 100
		bar := strings.Repeat("█", int(pct/5))
		fmt.Fprintf(&sb, "  overflow  | %-20s %6.1f%% (%d)\n", bar, pct, h.overflow)
	}
	return sb.String()
}
