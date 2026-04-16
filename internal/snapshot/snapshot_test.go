package snapshot_test

import (
	"testing"
	"time"

	"github.com/nicklaw5/grpcannon/internal/metrics"
	"github.com/nicklaw5/grpcannon/internal/snapshot"
)

func newRecorder() *metrics.Recorder {
	return metrics.NewRecorder()
}

func TestCapture_EmptyRecorder_ZeroCounts(t *testing.T) {
	r := newRecorder()
	c := snapshot.NewCollector(r)
	s := c.Capture()
	if s.Total != 0 {
		t.Fatalf("expected 0 total, got %d", s.Total)
	}
	if s.Errors != 0 {
		t.Fatalf("expected 0 errors, got %d", s.Errors)
	}
}

func TestCapture_RecordsElapsed(t *testing.T) {
	r := newRecorder()
	c := snapshot.NewCollector(r)
	time.Sleep(5 * time.Millisecond)
	s := c.Capture()
	if s.Elapsed < 5*time.Millisecond {
		t.Fatalf("expected elapsed >= 5ms, got %v", s.Elapsed)
	}
}

func TestAll_ReturnsAllSnapshots(t *testing.T) {
	r := newRecorder()
	c := snapshot.NewCollector(r)
	c.Capture()
	c.Capture()
	c.Capture()
	if len(c.All()) != 3 {
		t.Fatalf("expected 3 snapshots, got %d", len(c.All()))
	}
}

func TestLatest_NoSnapshots_ReturnsZero(t *testing.T) {
	r := newRecorder()
	c := snapshot.NewCollector(r)
	s := c.Latest()
	if s.Total != 0 || s.Elapsed != 0 {
		t.Fatal("expected zero snapshot")
	}
}

func TestLatest_ReturnsLastCaptured(t *testing.T) {
	r := newRecorder()
	c := snapshot.NewCollector(r)
	c.Capture()
	time.Sleep(2 * time.Millisecond)
	c.Capture()
	all := c.All()
	latest := c.Latest()
	if latest.Timestamp != all[len(all)-1].Timestamp {
		t.Fatal("latest does not match last snapshot")
	}
}

func TestCapture_RPSIsPositiveAfterDelay(t *testing.T) {
	r := newRecorder()
	r.Record(10*time.Millisecond, nil)
	c := snapshot.NewCollector(r)
	time.Sleep(5 * time.Millisecond)
	s := c.Capture()
	if s.RPS <= 0 {
		t.Fatalf("expected positive RPS, got %f", s.RPS)
	}
}
