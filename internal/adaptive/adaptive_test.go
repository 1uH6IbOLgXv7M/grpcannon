package adaptive_test

import (
	"testing"

	"github.com/nickpoorman/grpcannon/internal/adaptive"
)

func TestNew_DefaultCurrent_IsMin(t *testing.T) {
	c := adaptive.New(2, 10, 1)
	if got := c.Current(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestNew_ClampsBadMin(t *testing.T) {
	c := adaptive.New(0, 10, 1)
	if got := c.Current(); got < 1 {
		t.Fatalf("expected at least 1, got %d", got)
	}
}

func TestAdjust_NoRecords_ReturnsCurrent(t *testing.T) {
	c := adaptive.New(1, 10, 1)
	got := c.Adjust()
	if got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestAdjust_LowErrorRate_IncreasesConcurrency(t *testing.T) {
	c := adaptive.New(1, 10, 1)
	for i := 0; i < 100; i++ {
		c.Record(false)
	}
	got := c.Adjust()
	if got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestAdjust_HighErrorRate_ReducesConcurrency(t *testing.T) {
	c := adaptive.New(1, 10, 1)
	// Start at a higher level.
	c.Record(false)
	c.Adjust() // -> 2
	c.Record(false)
	c.Adjust() // -> 3

	// Now inject errors above threshold.
	for i := 0; i < 5; i++ {
		c.Record(true)
	}
	for i := 0; i < 5; i++ {
		c.Record(false)
	}
	got := c.Adjust()
	if got >= 3 {
		t.Fatalf("expected concurrency to drop below 3, got %d", got)
	}
}

func TestAdjust_ClampsAtMax(t *testing.T) {
	c := adaptive.New(1, 3, 1)
	for i := 0; i < 10; i++ {
		c.Record(false)
		c.Adjust()
	}
	if got := c.Current(); got > 3 {
		t.Fatalf("expected at most 3, got %d", got)
	}
}

func TestAdjust_ClampsAtMin(t *testing.T) {
	c := adaptive.New(2, 10, 1)
	for i := 0; i < 20; i++ {
		c.Record(true)
	}
	c.Adjust()
	if got := c.Current(); got < 2 {
		t.Fatalf("expected at least 2, got %d", got)
	}
}

func TestAdjust_ResetsCounters(t *testing.T) {
	c := adaptive.New(1, 10, 1)
	for i := 0; i < 10; i++ {
		c.Record(true)
	}
	c.Adjust()
	// Second adjust with no new records should not change current.
	before := c.Current()
	c.Adjust()
	if got := c.Current(); got != before {
		t.Fatalf("expected %d, got %d", before, got)
	}
}
