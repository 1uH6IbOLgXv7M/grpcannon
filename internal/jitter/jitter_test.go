package jitter

import (
	"testing"
	"time"
)

// fixed returns a Source that always returns v.
func fixed(v float64) Source { return func() float64 { return v } }

func TestFull_ZeroDuration_ReturnsZero(t *testing.T) {
	if got := full(0, fixed(0.5)); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestFull_ReturnsScaledDuration(t *testing.T) {
	d := 100 * time.Millisecond
	got := full(d, fixed(0.5))
	want := 50 * time.Millisecond
	if got != want {
		t.Fatalf("want %v got %v", want, got)
	}
}

func TestFull_MaxBound(t *testing.T) {
	d := 100 * time.Millisecond
	got := full(d, fixed(0.999))
	if got >= d {
		t.Fatalf("expected < %v, got %v", d, got)
	}
}

func TestEqual_ZeroDuration_ReturnsZero(t *testing.T) {
	if got := equal(0, fixed(0.5)); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestEqual_ReturnsInLowerHalfRange(t *testing.T) {
	d := 100 * time.Millisecond
	got := equal(d, fixed(0.0))
	want := 50 * time.Millisecond
	if got != want {
		t.Fatalf("want %v got %v", want, got)
	}
}

func TestEqual_ReturnsInUpperHalfRange(t *testing.T) {
	d := 100 * time.Millisecond
	got := equal(d, fixed(1.0))
	if got > d {
		t.Fatalf("expected <= %v, got %v", d, got)
	}
	if got < d/2 {
		t.Fatalf("expected >= %v, got %v", d/2, got)
	}
}

func TestDeviation_ZeroDuration_ReturnsZero(t *testing.T) {
	if got := deviation(0, 0.1, fixed(0.5)); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestDeviation_ZeroFactor_ReturnsBase(t *testing.T) {
	d := 100 * time.Millisecond
	if got := deviation(d, 0, fixed(0.5)); got != d {
		t.Fatalf("want %v got %v", d, got)
	}
}

func TestDeviation_FactorClamped(t *testing.T) {
	d := 100 * time.Millisecond
	// factor > 1 clamped to 1 → range [0, 200ms)
	got := deviation(d, 2.0, fixed(0.5))
	if got < 0 {
		t.Fatalf("negative duration: %v", got)
	}
}

func TestDeviation_MidPoint(t *testing.T) {
	d := 100 * time.Millisecond
	// factor=0.1 → delta=10ms, min=90ms; src=0.5 → 90+10=100ms
	got := deviation(d, 0.1, fixed(0.5))
	want := 100 * time.Millisecond
	if got != want {
		t.Fatalf("want %v got %v", want, got)
	}
}
