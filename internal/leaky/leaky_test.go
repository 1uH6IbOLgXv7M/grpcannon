package leaky_test

import (
	"context"
	"testing"
	"time"

	"github.com/your-org/grpcannon/internal/leaky"
)

func TestNew_ZeroCapacity_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero capacity")
		}
	}()
	leaky.New(0, 10)
}

func TestNew_ZeroRate_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero rate")
		}
	}()
	leaky.New(10, 0)
}

func TestAcquire_BelowCapacity_Succeeds(t *testing.T) {
	b := leaky.New(5, 1)
	for i := 0; i < 5; i++ {
		if err := b.Acquire(context.Background()); err != nil {
			t.Fatalf("unexpected error on attempt %d: %v", i, err)
		}
	}
}

func TestAcquire_AtCapacity_ReturnsErrDropped(t *testing.T) {
	b := leaky.New(3, 1)
	for i := 0; i < 3; i++ {
		_ = b.Acquire(context.Background())
	}
	if err := b.Acquire(context.Background()); err != leaky.ErrDropped {
		t.Fatalf("want ErrDropped, got %v", err)
	}
}

func TestAcquire_CancelledContext_ReturnsContextError(t *testing.T) {
	b := leaky.New(10, 5)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := b.Acquire(ctx); err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestAcquire_DrainOverTime_AllowsNewTokens(t *testing.T) {
	// capacity=2, rate=100/s — bucket drains in ~20 ms
	b := leaky.New(2, 100)
	_ = b.Acquire(context.Background())
	_ = b.Acquire(context.Background())

	// Bucket is full; wait for drain.
	time.Sleep(30 * time.Millisecond)

	if err := b.Acquire(context.Background()); err != nil {
		t.Fatalf("expected bucket to have drained, got %v", err)
	}
}

func TestInFlight_ReflectsCurrentLevel(t *testing.T) {
	b := leaky.New(10, 1)
	for i := 0; i < 4; i++ {
		_ = b.Acquire(context.Background())
	}
	if got := b.InFlight(); got < 3 || got > 5 {
		t.Fatalf("unexpected InFlight value: %f", got)
	}
}
