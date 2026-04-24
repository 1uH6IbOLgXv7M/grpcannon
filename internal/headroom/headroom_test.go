package headroom

import (
	"testing"
)

func TestNew_ClampsMaxConcurrency(t *testing.T) {
	e := New(Config{MaxConcurrency: 0})
	if e.cfg.MaxConcurrency != 1 {
		t.Fatalf("expected MaxConcurrency=1, got %d", e.cfg.MaxConcurrency)
	}
}

func TestNew_ClampsErrorRateThreshold(t *testing.T) {
	e := New(Config{MaxConcurrency: 10, ErrorRateThreshold: -0.5})
	if e.cfg.ErrorRateThreshold != 0 {
		t.Fatalf("expected threshold=0, got %f", e.cfg.ErrorRateThreshold)
	}
	e2 := New(Config{MaxConcurrency: 10, ErrorRateThreshold: 1.5})
	if e2.cfg.ErrorRateThreshold != 1 {
		t.Fatalf("expected threshold=1, got %f", e2.cfg.ErrorRateThreshold)
	}
}

func TestScore_NoLoad_ReturnsOne(t *testing.T) {
	e := New(Config{MaxConcurrency: 10})
	if got := e.Score(); got != 1.0 {
		t.Fatalf("expected 1.0, got %f", got)
	}
}

func TestScore_FullLoad_ReturnsZero(t *testing.T) {
	e := New(Config{MaxConcurrency: 10})
	e.Update(10, 0)
	if got := e.Score(); got != 0.0 {
		t.Fatalf("expected 0.0, got %f", got)
	}
}

func TestScore_HalfLoad_ReturnsHalf(t *testing.T) {
	e := New(Config{MaxConcurrency: 10})
	e.Update(5, 0)
	if got := e.Score(); got != 0.5 {
		t.Fatalf("expected 0.5, got %f", got)
	}
}

func TestScore_OverLoad_ClampsToZero(t *testing.T) {
	e := New(Config{MaxConcurrency: 10})
	e.Update(20, 0)
	if got := e.Score(); got != 0.0 {
		t.Fatalf("expected 0.0, got %f", got)
	}
}

func TestScore_ErrorRateExceedsThreshold_ReturnsZero(t *testing.T) {
	e := New(Config{MaxConcurrency: 10, ErrorRateThreshold: 0.5})
	e.Update(2, 0.6)
	if got := e.Score(); got != 0.0 {
		t.Fatalf("expected 0.0 due to error rate, got %f", got)
	}
}

func TestScore_ErrorRateBelowThreshold_UsesLoadRatio(t *testing.T) {
	e := New(Config{MaxConcurrency: 10, ErrorRateThreshold: 0.5})
	e.Update(2, 0.1)
	want := 0.8
	if got := e.Score(); got != want {
		t.Fatalf("expected %f, got %f", want, got)
	}
}

func TestAvailable_NoLoad_ReturnsTrue(t *testing.T) {
	e := New(Config{MaxConcurrency: 10})
	if !e.Available() {
		t.Fatal("expected Available=true")
	}
}

func TestAvailable_FullLoad_ReturnsFalse(t *testing.T) {
	e := New(Config{MaxConcurrency: 10})
	e.Update(10, 0)
	if e.Available() {
		t.Fatal("expected Available=false")
	}
}
