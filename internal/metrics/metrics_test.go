package metrics

import (
	"errors"
	"testing"
	"time"
)

func TestRecorder_EmptySnapshot(t *testing.T) {
	r := NewRecorder()
	s := r.Snapshot()
	if s.Total != 0 {
		t.Errorf("expected Total 0, got %d", s.Total)
	}
	if s.Min != 0 || s.Max != 0 {
		t.Errorf("expected zero durations for empty recorder")
	}
}

func TestRecorder_CountsErrors(t *testing.T) {
	r := NewRecorder()
	r.Record(10*time.Millisecond, nil)
	r.Record(20*time.Millisecond, errors.New("rpc error"))
	r.Record(30*time.Millisecond, nil)

	s := r.Snapshot()
	if s.Total != 3 {
		t.Errorf("expected Total 3, got %d", s.Total)
	}
	if s.Errors != 1 {
		t.Errorf("expected Errors 1, got %d", s.Errors)
	}
}

func TestRecorder_Percentiles(t *testing.T) {
	r := NewRecorder()
	for i := 1; i <= 100; i++ {
		r.Record(time.Duration(i)*time.Millisecond, nil)
	}

	s := r.Snapshot()

	if s.Min != 1*time.Millisecond {
		t.Errorf("expected Min 1ms, got %v", s.Min)
	}
	if s.Max != 100*time.Millisecond {
		t.Errorf("expected Max 100ms, got %v", s.Max)
	}
	if s.P50 != 50*time.Millisecond {
		t.Errorf("expected P50 50ms, got %v", s.P50)
	}
	if s.P95 != 95*time.Millisecond {
		t.Errorf("expected P95 95ms, got %v", s.P95)
	}
	if s.P99 != 99*time.Millisecond {
		t.Errorf("expected P99 99ms, got %v", s.P99)
	}
}

func TestRecorder_Mean(t *testing.T) {
	r := NewRecorder()
	r.Record(10*time.Millisecond, nil)
	r.Record(20*time.Millisecond, nil)
	r.Record(30*time.Millisecond, nil)

	s := r.Snapshot()
	if s.Mean != 20*time.Millisecond {
		t.Errorf("expected Mean 20ms, got %v", s.Mean)
	}
}

func TestSummary_String(t *testing.T) {
	s := Summary{
		Total:  10,
		Errors: 2,
		Min:    5 * time.Millisecond,
		Max:    50 * time.Millisecond,
		Mean:   25 * time.Millisecond,
		P50:    24 * time.Millisecond,
		P95:    48 * time.Millisecond,
		P99:    50 * time.Millisecond,
	}
	out := s.String()
	if len(out) == 0 {
		t.Error("expected non-empty string from Summary.String()")
	}
}

func TestRecorder_ConcurrentWrites(t *testing.T) {
	r := NewRecorder()
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func(n int) {
			r.Record(time.Duration(n)*time.Millisecond, nil)
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 50; i++ {
		<-done
	}
	s := r.Snapshot()
	if s.Total != 50 {
		t.Errorf("expected Total 50, got %d", s.Total)
	}
}
