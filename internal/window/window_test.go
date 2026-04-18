package window

import (
	"testing"
	"time"
)

func TestCounts_Empty(t *testing.T) {
	w := New(time.Second, 10)
	total, errors := w.Counts()
	if total != 0 || errors != 0 {
		t.Fatalf("expected 0,0 got %d,%d", total, errors)
	}
}

func TestAdd_RecordsTotals(t *testing.T) {
	w := New(time.Second, 10)
	w.Add(false)
	w.Add(false)
	w.Add(true)
	total, errors := w.Counts()
	if total != 3 {
		t.Fatalf("expected total 3, got %d", total)
	}
	if errors != 1 {
		t.Fatalf("expected errors 1, got %d", errors)
	}
}

func TestErrorRate_NoRequests(t *testing.T) {
	w := New(time.Second, 10)
	if r := w.ErrorRate(); r != 0 {
		t.Fatalf("expected 0, got %f", r)
	}
}

func TestErrorRate_AllErrors(t *testing.T) {
	w := New(time.Second, 10)
	w.Add(true)
	w.Add(true)
	if r := w.ErrorRate(); r != 1.0 {
		t.Fatalf("expected 1.0, got %f", r)
	}
}

func TestErrorRate_HalfErrors(t *testing.T) {
	w := New(time.Second, 10)
	w.Add(false)
	w.Add(true)
	if r := w.ErrorRate(); r != 0.5 {
		t.Fatalf("expected 0.5, got %f", r)
	}
}

func TestCounts_ExpiredSlots_NotCounted(t *testing.T) {
	now := time.Unix(1000, 0)
	w := New(100*time.Millisecond, 10)
	w.now = func() time.Time { return now }
	w.Add(false)
	w.Add(true)
	// advance past full window
	w.now = func() time.Time { return now.Add(200 * time.Millisecond) }
	total, errors := w.Counts()
	if total != 0 || errors != 0 {
		t.Fatalf("expected expired slots to be ignored, got %d,%d", total, errors)
	}
}

func TestNew_SingleSlot(t *testing.T) {
	w := New(time.Second, 0) // should clamp to 1
	if w.size != 1 {
		t.Fatalf("expected size 1, got %d", w.size)
	}
}
