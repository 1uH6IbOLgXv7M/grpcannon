package overflow_test

import (
	"sync"
	"testing"

	"github.com/example/grpcannon/internal/overflow"
)

func TestNew_CapacityLessThanOne_ClampsToOne(t *testing.T) {
	q := overflow.New(0)
	if !q.Acquire() {
		t.Fatal("expected first acquire to succeed")
	}
	if q.Acquire() {
		t.Fatal("expected second acquire to be dropped")
	}
}

func TestAcquire_BelowCapacity_Succeeds(t *testing.T) {
	q := overflow.New(3)
	for i := 0; i < 3; i++ {
		if !q.Acquire() {
			t.Fatalf("acquire %d should succeed", i)
		}
	}
}

func TestAcquire_AtCapacity_Drops(t *testing.T) {
	q := overflow.New(2)
	q.Acquire()
	q.Acquire()
	if q.Acquire() {
		t.Fatal("expected drop when at capacity")
	}
	if q.Dropped() != 1 {
		t.Fatalf("expected 1 dropped, got %d", q.Dropped())
	}
}

func TestRelease_FreesSlot(t *testing.T) {
	q := overflow.New(1)
	q.Acquire()
	q.Release()
	if !q.Acquire() {
		t.Fatal("expected acquire to succeed after release")
	}
}

func TestTotal_TracksAllAttempts(t *testing.T) {
	q := overflow.New(2)
	for i := 0; i < 5; i++ {
		q.Acquire()
	}
	if q.Total() != 5 {
		t.Fatalf("expected total 5, got %d", q.Total())
	}
}

func TestInFlight_ReflectsAcquired(t *testing.T) {
	q := overflow.New(4)
	q.Acquire()
	q.Acquire()
	if q.InFlight() != 2 {
		t.Fatalf("expected in-flight 2, got %d", q.InFlight())
	}
}

func TestConcurrent_NeverExceedsCapacity(t *testing.T) {
	const cap = 10
	q := overflow.New(cap)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if q.Acquire() {
				if q.InFlight() > cap {
					t.Errorf("in-flight exceeded capacity")
				}
				q.Release()
			}
		}()
	}
	wg.Wait()
}
