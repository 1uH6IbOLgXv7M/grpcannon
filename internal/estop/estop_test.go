package estop

import (
	"sync"
	"testing"
)

func TestAllow_NotTripped_ReturnsNil(t *testing.T) {
	e := New(0.5)
	if err := e.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecordFailure_BelowThreshold_DoesNotTrip(t *testing.T) {
	e := New(0.5)
	e.RecordSuccess()
	e.RecordSuccess()
	e.RecordFailure() // 1/3 ≈ 0.33 < 0.5
	if e.Tripped() {
		t.Fatal("expected not tripped")
	}
	if err := e.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecordFailure_ExceedsThreshold_Trips(t *testing.T) {
	e := New(0.5)
	e.RecordFailure()
	e.RecordFailure() // 2/2 = 1.0 > 0.5
	if !e.Tripped() {
		t.Fatal("expected tripped")
	}
	if err := e.Allow(); err != ErrTripped {
		t.Fatalf("expected ErrTripped, got %v", err)
	}
}

func TestReset_ClearsTrippedState(t *testing.T) {
	e := New(0.1)
	e.RecordFailure()
	if !e.Tripped() {
		t.Fatal("expected tripped before reset")
	}
	e.Reset()
	if e.Tripped() {
		t.Fatal("expected not tripped after reset")
	}
	if err := e.Allow(); err != nil {
		t.Fatalf("expected nil after reset, got %v", err)
	}
}

func TestReset_ClearsCounters(t *testing.T) {
	e := New(0.5)
	e.RecordFailure()
	e.RecordFailure()
	e.Reset()
	if got := e.ErrorRate(); got != 0 {
		t.Fatalf("expected 0 error rate after reset, got %f", got)
	}
}

func TestErrorRate_NoObservations_ReturnsZero(t *testing.T) {
	e := New(0.5)
	if r := e.ErrorRate(); r != 0 {
		t.Fatalf("expected 0, got %f", r)
	}
}

func TestErrorRate_HalfErrors(t *testing.T) {
	e := New(0.9)
	e.RecordSuccess()
	e.RecordFailure()
	if got, want := e.ErrorRate(), 0.5; got != want {
		t.Fatalf("expected %f, got %f", want, got)
	}
}

func TestNew_ClampsThreshold(t *testing.T) {
	if e := New(-1); e.threshold != 0 {
		t.Fatalf("expected threshold clamped to 0, got %f", e.threshold)
	}
	if e := New(2); e.threshold != 1 {
		t.Fatalf("expected threshold clamped to 1, got %f", e.threshold)
	}
}

func TestConcurrent_RecordAndAllow(t *testing.T) {
	e := New(0.9)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				e.RecordSuccess()
			} else {
				e.RecordFailure()
			}
			_ = e.Allow()
		}(i)
	}
	wg.Wait()
	// 50/100 = 0.5 < 0.9, should not be tripped
	if e.Tripped() {
		t.Fatal("expected not tripped at 50% error rate with 0.9 threshold")
	}
}
