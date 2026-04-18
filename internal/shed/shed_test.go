package shed

import (
	"sync"
	"testing"
)

func TestAcquire_Unlimited_NeverSheds(t *testing.T) {
	s := New(0)
	for i := 0; i < 1000; i++ {
		if err := s.Acquire(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}

func TestAcquire_AtLimit_ReturnsErrShed(t *testing.T) {
	s := New(2)
	if err := s.Acquire(); err != nil {
		t.Fatal(err)
	}
	if err := s.Acquire(); err != nil {
		t.Fatal(err)
	}
	if err := s.Acquire(); err != ErrShed {
		t.Fatalf("expected ErrShed, got %v", err)
	}
}

func TestRelease_FreesSlot(t *testing.T) {
	s := New(1)
	if err := s.Acquire(); err != nil {
		t.Fatal(err)
	}
	s.Release()
	if err := s.Acquire(); err != nil {
		t.Fatalf("expected slot after release, got %v", err)
	}
}

func TestInFlight_TracksCount(t *testing.T) {
	s := New(10)
	s.Acquire()
	s.Acquire()
	if got := s.InFlight(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
	s.Release()
	if got := s.InFlight(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestConcurrent_NeverExceedsLimit(t *testing.T) {
	const limit = 10
	const goroutines = 200
	s := New(limit)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.Acquire(); err == nil {
				if v := s.InFlight(); v > limit {
					t.Errorf("in-flight %d exceeded limit %d", v, limit)
				}
				s.Release()
			}
		}()
	}
	wg.Wait()
}
