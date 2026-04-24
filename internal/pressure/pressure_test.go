package pressure

import (
	"testing"
	"time"
)

func fixed(ts time.Time) func() time.Time {
	return func() time.Time { return ts }
}

func TestScore_NoObservations_ReturnsZero(t *testing.T) {
	tr := New(Config{})
	if got := tr.Score(); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestScore_SingleSuccessLowLatency_LowPressure(t *testing.T) {
	tr := New(Config{HighLatency: 2 * time.Second})
	tr.Record(100*time.Millisecond, false)
	got := tr.Score()
	if got >= 0.1 {
		t.Fatalf("expected low pressure, got %v", got)
	}
}

func TestScore_HighLatency_PressureApproachesOne(t *testing.T) {
	tr := New(Config{HighLatency: 1 * time.Second})
	tr.Record(1*time.Second, false)
	got := tr.Score()
	// latency component = 1.0, error component = 0 → 0.7
	if got < 0.69 || got > 0.71 {
		t.Fatalf("expected ~0.70, got %v", got)
	}
}

func TestScore_AllErrors_PressureIncludesErrorWeight(t *testing.T) {
	tr := New(Config{HighLatency: 2 * time.Second})
	tr.Record(0, true)
	tr.Record(0, true)
	got := tr.Score()
	// latency = 0, error rate = 1.0 → 0.3
	if got < 0.29 || got > 0.31 {
		t.Fatalf("expected ~0.30, got %v", got)
	}
}

func TestScore_FullSaturation_CapsAtOne(t *testing.T) {
	tr := New(Config{HighLatency: 500 * time.Millisecond})
	tr.Record(2*time.Second, true)
	got := tr.Score()
	if got > 1.0 {
		t.Fatalf("score must not exceed 1.0, got %v", got)
	}
	if got != 1.0 {
		t.Fatalf("expected 1.0, got %v", got)
	}
}

func TestEvict_RemovesStaleObservations(t *testing.T) {
	now := time.Now()
	tr := New(Config{Window: 5 * time.Second, HighLatency: 2 * time.Second})

	// inject an observation that is 10 s old
	tr.mu.Lock()
	tr.obs = append(tr.obs, observation{
		at:      now.Add(-10 * time.Second),
		latency: 2 * time.Second,
		err:     true,
	})
	tr.mu.Unlock()

	// Score should evict the stale entry and return 0
	if got := tr.Score(); got != 0 {
		t.Fatalf("expected 0 after eviction, got %v", got)
	}
}

func TestRecord_ConcurrentAccess_NoPanic(t *testing.T) {
	tr := New(Config{})
	done := make(chan struct{})
	for i := 0; i < 8; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				tr.Record(10*time.Millisecond, j%5 == 0)
				_ = tr.Score()
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < 8; i++ {
		<-done
	}
}
