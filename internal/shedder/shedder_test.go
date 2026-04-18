package shedder_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/nickbadlose/grpcannon/internal/shedder"
)

func TestAcquire_BelowLimit_Succeeds(t *testing.T) {
	s := shedder.New(shedder.Config{MaxInFlight: 5})
	if err := s.Acquire(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAcquire_AtLimit_ReturnsShed(t *testing.T) {
	s := shedder.New(shedder.Config{MaxInFlight: 1})
	if err := s.Acquire(); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	if err := s.Acquire(); !errors.Is(err, shedder.ErrShed) {
		t.Fatalf("expected ErrShed, got %v", err)
	}
}

func TestRelease_FreesSlot(t *testing.T) {
	s := shedder.New(shedder.Config{MaxInFlight: 1})
	_ = s.Acquire()
	s.Release(nil)
	if err := s.Acquire(); err != nil {
		t.Fatalf("expected slot free after release: %v", err)
	}
}

func TestAcquire_HighErrorRate_ReturnsShed(t *testing.T) {
	s := shedder.New(shedder.Config{
		ErrorRateThreshold: 0.5,
		WindowSize:         time.Second,
	})
	errSentinel := errors.New("rpc error")
	// Drive error rate above threshold.
	for i := 0; i < 10; i++ {
		_ = s.Acquire()
		s.Release(errSentinel)
	}
	if err := s.Acquire(); !errors.Is(err, shedder.ErrShed) {
		t.Fatalf("expected ErrShed due to high error rate, got %v", err)
	}
}

func TestInFlight_TracksCount(t *testing.T) {
	s := shedder.New(shedder.Config{MaxInFlight: 10})
	for i := 0; i < 3; i++ {
		_ = s.Acquire()
	}
	if got := s.InFlight(); got != 3 {
		t.Fatalf("expected 3 in-flight, got %d", got)
	}
}

func TestConcurrent_NeverExceedsLimit(t *testing.T) {
	const limit = 10
	s := shedder.New(shedder.Config{MaxInFlight: limit})
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.Acquire(); err == nil {
				if s.InFlight() > limit {
					t.Errorf("in-flight %d exceeded limit %d", s.InFlight(), limit)
				}
				s.Release(nil)
			}
		}()
	}
	wg.Wait()
}
