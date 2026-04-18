package shedder_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nickbadlose/grpcannon/internal/shedder"
)

// TestIntegration_ShedUnderLoad verifies that the shedder protects a simulated
// backend whose capacity is lower than the offered load.
func TestIntegration_ShedUnderLoad(t *testing.T) {
	const capacity = 5
	s := shedder.New(shedder.Config{
		MaxInFlight:        capacity,
		ErrorRateThreshold: 0.8,
		WindowSize:         2 * time.Second,
	})

	var (
		accepted atomic.Int64
		shed     atomic.Int64
		wg       sync.WaitGroup
	)

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.Acquire(); err != nil {
				if errors.Is(err, shedder.ErrShed) {
					shed.Add(1)
				}
				return
			}
			accepted.Add(1)
			time.Sleep(5 * time.Millisecond)
			s.Release(nil)
		}()
	}
	wg.Wait()

	if accepted.Load() == 0 {
		t.Fatal("expected some requests to be accepted")
	}
	if shed.Load() == 0 {
		t.Fatal("expected some requests to be shed under load")
	}
	t.Logf("accepted=%d shed=%d", accepted.Load(), shed.Load())
}
