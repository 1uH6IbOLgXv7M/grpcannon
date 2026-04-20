package eventlog

import (
	"testing"
)

func TestNew_CapacityLessThanOne_ClampsToOne(t *testing.T) {
	l := New(0)
	if l.capacity != 1 {
		t.Fatalf("expected capacity 1, got %d", l.capacity)
	}
}

func TestAdd_SingleEntry_RetainedInEntries(t *testing.T) {
	l := New(10)
	l.Add(LevelInfo, "hello")
	entries := l.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Message != "hello" {
		t.Errorf("unexpected message: %s", entries[0].Message)
	}
	if entries[0].Level != LevelInfo {
		t.Errorf("unexpected level: %v", entries[0].Level)
	}
}

func TestAdd_BelowCapacity_AllRetained(t *testing.T) {
	l := New(5)
	msgs := []string{"a", "b", "c"}
	for _, m := range msgs {
		l.Add(LevelWarn, m)
	}
	entries := l.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i, e := range entries {
		if e.Message != msgs[i] {
			t.Errorf("entry %d: expected %q got %q", i, msgs[i], e.Message)
		}
	}
}

func TestAdd_ExceedsCapacity_OldestOverwritten(t *testing.T) {
	l := New(3)
	for _, m := range []string{"a", "b", "c", "d"} {
		l.Add(LevelError, m)
	}
	entries := l.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	want := []string{"b", "c", "d"}
	for i, e := range entries {
		if e.Message != want[i] {
			t.Errorf("entry %d: expected %q got %q", i, want[i], e.Message)
		}
	}
}

func TestTotal_TracksAllAdded(t *testing.T) {
	l := New(2)
	for i := 0; i < 5; i++ {
		l.Add(LevelInfo, "x")
	}
	if l.Total() != 5 {
		t.Errorf("expected total 5, got %d", l.Total())
	}
}

func TestEntries_Empty_ReturnsEmptySlice(t *testing.T) {
	l := New(4)
	if entries := l.Entries(); len(entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(entries))
	}
}

func TestEntries_TimestampSet(t *testing.T) {
	l := New(1)
	l.Add(LevelInfo, "ts-check")
	e := l.Entries()[0]
	if e.At.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
