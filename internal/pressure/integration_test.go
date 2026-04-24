package pressure_test

import (
	"testing"
	"time"

	"github.com/example/grpcannon/internal/pressure"
)

// TestPressureRises_ThenFallsAfterWindow verifies that pressure decays
// naturally once observations age out of the retention window.
func TestPressureRises_ThenFallsAfterWindow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	cfg := pressure.Config{
		Window:      200 * time.Millisecond,
		HighLatency: 100 * time.Millisecond,
	}
	tr := pressure.New(cfg)

	// Drive pressure up.
	for i := 0; i < 20; i++ {
		tr.Record(150*time.Millisecond, true)
	}
	high := tr.Score()
	if high < 0.5 {
		t.Fatalf("expected high pressure, got %v", high)
	}

	// Wait for the window to expire.
	time.Sleep(250 * time.Millisecond)

	low := tr.Score()
	if low != 0 {
		t.Fatalf("expected pressure to drop to 0 after window, got %v", low)
	}
}
