package pacemaker_test

import (
	"sync"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/pacemaker"
)

// TestConcurrent_AdjustAndCurrent verifies that concurrent calls to Adjust and
// Current do not race (run with -race).
func TestConcurrent_AdjustAndCurrent(t *testing.T) {
	pm := pacemaker.New(pacemaker.Config{
		TargetP99:  50 * time.Millisecond,
		MinRPS:     1,
		MaxRPS:     1000,
		StepFactor: 0.05,
	})

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			latency := time.Duration(i*5) * time.Millisecond
			for j := 0; j < 50; j++ {
				pm.Adjust(latency)
				_ = pm.Current()
			}
		}(i)
	}
	wg.Wait()

	v := pm.Current()
	if v < 1 || v > 1000 {
		t.Fatalf("rate out of bounds: %f", v)
	}
}
