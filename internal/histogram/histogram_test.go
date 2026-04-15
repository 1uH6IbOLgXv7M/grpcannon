package histogram

import (
	"strings"
	"testing"
	"time"
)

func TestNew_EmptyHistogram_TotalIsZero(t *testing.T) {
	h := New()
	if got := h.Total(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestObserve_SingleObservation_TotalIsOne(t *testing.T) {
	h := New()
	h.Observe(5 * time.Millisecond)
	if got := h.Total(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestObserve_FallsIntoCorrectBucket(t *testing.T) {
	h := NewWithBuckets([]float64{10, 50, 100})
	h.Observe(8 * time.Millisecond) // bucket 0 (<= 10 ms)
	h.Observe(45 * time.Millisecond) // bucket 1 (<= 50 ms)
	h.Observe(99 * time.Millisecond) // bucket 2 (<= 100 ms)

	if h.counts[0] != 1 {
		t.Errorf("bucket 0: expected 1, got %d", h.counts[0])
	}
	if h.counts[1] != 1 {
		t.Errorf("bucket 1: expected 1, got %d", h.counts[1])
	}
	if h.counts[2] != 1 {
		t.Errorf("bucket 2: expected 1, got %d", h.counts[2])
	}
}

func TestObserve_OverflowBucket(t *testing.T) {
	h := NewWithBuckets([]float64{1, 5})
	h.Observe(999 * time.Millisecond)
	if h.overflow != 1 {
		t.Fatalf("expected overflow=1, got %d", h.overflow)
	}
	if h.Total() != 1 {
		t.Fatalf("total should include overflow, got %d", h.Total())
	}
}

func TestObserve_MultipleObservations_TotalMatches(t *testing.T) {
	h := New()
	durations := []time.Duration{1, 2, 5, 10, 50, 100, 500, 1000}
	for _, d := range durations {
		h.Observe(d * time.Millisecond)
	}
	if got := h.Total(); got != int64(len(durations)) {
		t.Fatalf("expected %d, got %d", len(durations), got)
	}
}

func TestString_EmptyHistogram(t *testing.T) {
	h := New()
	got := h.String()
	if got != "(no observations)" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestString_ContainsBucketLabels(t *testing.T) {
	h := New()
	h.Observe(3 * time.Millisecond)
	out := h.String()
	if !strings.Contains(out, "ms") {
		t.Errorf("expected 'ms' in output, got:\n%s", out)
	}
}

func TestString_OverflowLinePresent(t *testing.T) {
	h := NewWithBuckets([]float64{1})
	h.Observe(500 * time.Millisecond)
	out := h.String()
	if !strings.Contains(out, "overflow") {
		t.Errorf("expected overflow line, got:\n%s", out)
	}
}
