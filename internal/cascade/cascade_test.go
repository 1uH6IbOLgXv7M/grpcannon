package cascade_test

import (
	"testing"

	"github.com/your-org/grpcannon/internal/cascade"
)

func TestNew_ClampsBadThreshold(t *testing.T) {
	d := cascade.New(0)
	d.RecordFailure()
	if !d.Tripped() {
		t.Fatal("expected detector to trip after one failure when threshold clamped to 1")
	}
}

func TestAllow_ClosedByDefault(t *testing.T) {
	d := cascade.New(3)
	if err := d.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecordFailure_TripsAfterThreshold(t *testing.T) {
	d := cascade.New(3)
	d.RecordFailure()
	d.RecordFailure()
	if d.Tripped() {
		t.Fatal("should not be tripped before threshold")
	}
	d.RecordFailure()
	if !d.Tripped() {
		t.Fatal("expected tripped after threshold")
	}
	if err := d.Allow(); err == nil {
		t.Fatal("expected ErrOpen")
	}
}

func TestRecordSuccess_ResetsState(t *testing.T) {
	d := cascade.New(2)
	d.RecordFailure()
	d.RecordFailure()
	if !d.Tripped() {
		t.Fatal("expected tripped")
	}
	d.RecordSuccess()
	if d.Tripped() {
		t.Fatal("expected reset after success")
	}
	if d.Consecutive() != 0 {
		t.Fatalf("expected consecutive=0, got %d", d.Consecutive())
	}
	if err := d.Allow(); err != nil {
		t.Fatalf("expected nil after reset, got %v", err)
	}
}

func TestConsecutive_TracksCount(t *testing.T) {
	d := cascade.New(10)
	for i := 1; i <= 5; i++ {
		d.RecordFailure()
		if got := d.Consecutive(); got != i {
			t.Fatalf("step %d: expected %d, got %d", i, i, got)
		}
	}
}

func TestConcurrent_RecordAndAllow(t *testing.T) {
	d := cascade.New(50)
	done := make(chan struct{})
	go func() {
		for i := 0; i < 200; i++ {
			d.RecordFailure()
		}
		close(done)
	}()
	for {
		select {
		case <-done:
			return
		default:
			_ = d.Allow()
			_ = d.Consecutive()
		}
	}
}
