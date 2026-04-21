package counter_test

import (
	"sync"
	"testing"

	"github.com/your-org/grpcannon/internal/counter"
)

func TestNew_ZeroValues(t *testing.T) {
	c := counter.New()
	if c.Total() != 0 {
		t.Fatalf("expected Total=0, got %d", c.Total())
	}
	if c.Errors() != 0 {
		t.Fatalf("expected Errors=0, got %d", c.Errors())
	}
}

func TestIncTotal_IncreasesTotal(t *testing.T) {
	c := counter.New()
	c.IncTotal()
	c.IncTotal()
	if c.Total() != 2 {
		t.Fatalf("expected Total=2, got %d", c.Total())
	}
	if c.Errors() != 0 {
		t.Fatalf("expected Errors=0, got %d", c.Errors())
	}
}

func TestIncErrors_IncreasesTotalAndErrors(t *testing.T) {
	c := counter.New()
	c.IncErrors()
	if c.Total() != 1 {
		t.Fatalf("expected Total=1, got %d", c.Total())
	}
	if c.Errors() != 1 {
		t.Fatalf("expected Errors=1, got %d", c.Errors())
	}
}

func TestErrorRate_NoRequests_ReturnsZero(t *testing.T) {
	c := counter.New()
	if r := c.ErrorRate(); r != 0 {
		t.Fatalf("expected 0, got %f", r)
	}
}

func TestErrorRate_HalfErrors(t *testing.T) {
	c := counter.New()
	c.IncTotal()
	c.IncErrors()
	if r := c.ErrorRate(); r != 0.5 {
		t.Fatalf("expected 0.5, got %f", r)
	}
}

func TestReset_ZeroesBothCounters(t *testing.T) {
	c := counter.New()
	c.IncTotal()
	c.IncErrors()
	c.Reset()
	if c.Total() != 0 || c.Errors() != 0 {
		t.Fatalf("expected both counters to be 0 after Reset")
	}
}

func TestSnapshot_CapturesCurrentState(t *testing.T) {
	c := counter.New()
	c.IncTotal()
	c.IncTotal()
	c.IncErrors()
	s := cif s.Total != 3 {
		t.Fatalf("expected Snapshot.Total=3, got %d", s.Total)
	}
	if s.Errors != 1 {
		t.Fatalf("expected Snapshot.Errors=1, got %d", s.Errors)
	}
}

func TestConcurrent_NeverRaces(t *testing.T) {
	c := counter.New()
	var wg sync.WaitGroup
	const goroutines = 50
	wg.Add(goroutines * 2)
	for i := 0; i < goroutines; i++ {
		go func() { defer wg.Done(); c.IncTotal() }()
		go func() { defer wg.Done(); c.IncErrors() }()
	}
	wg.Wait()
	if c.Total() != goroutines*2 {
		t.Fatalf("expected Total=%d, got %d", goroutines*2, c.Total())
	}
	if c.Errors() != goroutines {
		t.Fatalf("expected Errors=%d, got %d", goroutines, c.Errors())
	}
}
