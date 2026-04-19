package quorum

import (
	"errors"
	"testing"
)

func TestNew_ClampsMinRatio(t *testing.T) {
	q := New(-0.5, 1)
	if q.minRatio != 0 {
		t.Fatalf("expected 0, got %f", q.minRatio)
	}
	q2 := New(1.5, 1)
	if q2.minRatio != 1 {
		t.Fatalf("expected 1, got %f", q2.minRatio)
	}
}

func TestNew_ClampsMinTotal(t *testing.T) {
	q := New(0.9, 0)
	if q.minTotal != 1 {
		t.Fatalf("expected 1, got %d", q.minTotal)
	}
}

func TestRatio_NoObservations_ReturnsOne(t *testing.T) {
	q := New(0.9, 5)
	if q.Ratio() != 1.0 {
		t.Fatalf("expected 1.0, got %f", q.Ratio())
	}
}

func TestCheck_BelowMinTotal_ReturnsNil(t *testing.T) {
	q := New(0.9, 10)
	for i := 0; i < 5; i++ {
		q.RecordFailure()
	}
	if err := q.Check(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCheck_AboveThreshold_ReturnsNil(t *testing.T) {
	q := New(0.8, 5)
	for i := 0; i < 9; i++ {
		q.RecordSuccess()
	}
	q.RecordFailure()
	if err := q.Check(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCheck_BelowThreshold_ReturnsError(t *testing.T) {
	q := New(0.9, 5)
	for i := 0; i < 5; i++ {
		q.RecordSuccess()
	}
	for i := 0; i < 5; i++ {
		q.RecordFailure()
	}
	err := q.Check()
	if !errors.Is(err, ErrBelowThreshold) {
		t.Fatalf("expected ErrBelowThreshold, got %v", err)
	}
}

func TestTotal_TracksAllObservations(t *testing.T) {
	q := New(0.5, 1)
	q.RecordSuccess()
	q.RecordSuccess()
	q.RecordFailure()
	if q.Total() != 3 {
		t.Fatalf("expected 3, got %d", q.Total())
	}
}
