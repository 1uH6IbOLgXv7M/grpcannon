package budget_test

import (
	"errors"
	"testing"

	"github.com/your-org/grpcannon/internal/budget"
)

var errFake = errors.New("fake")

func TestAllow_NoRecords_ReturnsNil(t *testing.T) {
	b := budget.New(0.5)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllow_BelowThreshold_ReturnsNil(t *testing.T) {
	b := budget.New(0.5)
	b.Record(nil)
	b.Record(nil)
	b.Record(errFake)
	// ratio = 0.33 < 0.5
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllow_AtThreshold_ReturnsExceeded(t *testing.T) {
	b := budget.New(0.5)
	b.Record(nil)
	b.Record(errFake)
	// ratio = 0.5 >= 0.5
	if err := b.Allow(); !errors.Is(err, budget.ErrExceeded) {
		t.Fatalf("expected ErrExceeded, got %v", err)
	}
}

func TestAllow_AboveThreshold_ReturnsExceeded(t *testing.T) {
	b := budget.New(0.2)
	for i := 0; i < 3; i++ {
		b.Record(errFake)
	}
	b.Record(nil)
	if err := b.Allow(); !errors.Is(err, budget.ErrExceeded) {
		t.Fatalf("expected ErrExceeded, got %v", err)
	}
}

func TestRatio_AllFailures(t *testing.T) {
	b := budget.New(0.5)
	b.Record(errFake)
	b.Record(errFake)
	if got := b.Ratio(); got != 1.0 {
		t.Fatalf("expected 1.0, got %f", got)
	}
}

func TestReset_ClearsCounters(t *testing.T) {
	b := budget.New(0.5)
	b.Record(errFake)
	b.Record(errFake)
	b.Reset()
	if r := b.Ratio(); r != 0 {
		t.Fatalf("expected 0 after reset, got %f", r)
	}
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil after reset, got %v", err)
	}
}

func TestNew_ClampThreshold(t *testing.T) {
	// Should not panic or divide by zero with edge thresholds.
	b := budget.New(0)
	b.Record(nil)
	if b.Ratio() != 0 {
		t.Fatal("unexpected ratio")
	}
}
