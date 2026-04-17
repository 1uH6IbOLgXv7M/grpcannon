package circuit

import (
	"testing"
	"time"
)

func TestAllow_ClosedByDefault(t *testing.T) {
	b := New(3, time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecordFailure_OpensAfterThreshold(t *testing.T) {
	b := New(3, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	if b.CurrentState() != StateClosed {
		t.Fatal("expected closed after 2 failures")
	}
	b.RecordFailure()
	if b.CurrentState() != StateOpen {
		t.Fatal("expected open after 3 failures")
	}
}

func TestAllow_ReturnsErrOpenWhenOpen(t *testing.T) {
	b := New(1, time.Hour)
	b.RecordFailure()
	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestAllow_ClosesAfterCooldown(t *testing.T) {
	now := time.Now()
	b := New(1, 100*time.Millisecond)
	b.now = func() time.Time { return now }
	b.RecordFailure()

	// still open
	if err := b.Allow(); err != ErrOpen {
		t.Fatal("expected open")
	}

	// advance past cooldown
	b.now = func() time.Time { return now.Add(200 * time.Millisecond) }
	if err := b.Allow(); err != nil {
		t.Fatalf("expected closed after cooldown, got %v", err)
	}
	if b.CurrentState() != StateClosed {
		t.Fatal("expected state closed")
	}
}

func TestRecordSuccess_ResetsFailures(t *testing.T) {
	b := New(3, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	b.RecordSuccess()
	b.RecordFailure()
	b.RecordFailure()
	if b.CurrentState() != StateClosed {
		t.Fatal("expected closed: success should have reset counter")
	}
}

func TestNew_MinFailuresIsOne(t *testing.T) {
	b := New(0, time.Second)
	b.RecordFailure()
	if b.CurrentState() != StateOpen {
		t.Fatal("expected open with effective maxFailures=1")
	}
}

func TestRecordSuccess_ClosesOpenCircuit(t *testing.T) {
	b := New(1, time.Hour)
	b.RecordFailure()
	if b.CurrentState() != StateOpen {
		t.Fatal("expected open after failure")
	}
	b.RecordSuccess()
	if b.CurrentState() != StateClosed {
		t.Fatal("expected closed after success")
	}
	if err := b.Allow(); err != nil {
		t.Fatalf("expected allow after success reset, got %v", err)
	}
}
