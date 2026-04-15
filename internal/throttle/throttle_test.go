package throttle_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/throttle"
)

func TestNew_Unlimited_WaitReturnsImmediately(t *testing.T) {
	th := throttle.New(0)
	defer th.Stop()

	ctx := context.Background()
	done := make(chan error, 1)
	go func() { done <- th.Wait(ctx) }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatal("Wait blocked unexpectedly for unlimited throttle")
	}
}

func TestNew_RateLimited_ConsumesTokens(t *testing.T) {
	th := throttle.New(100)
	defer th.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for i := 0; i < 5; i++ {
		if err := th.Wait(ctx); err != nil {
			t.Fatalf("Wait[%d] returned unexpected error: %v", i, err)
		}
	}
}

func TestWait_ContextCancelled_ReturnsError(t *testing.T) {
	th := throttle.New(1) // very slow — 1 token/s
	defer th.Stop()

	// drain the first token so the next Wait will block
	ctx := context.Background()
	_ = th.Wait(ctx)

	ctxCancel, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	if err := th.Wait(ctxCancel); err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}

func TestStop_CalledTwice_NoPanic(t *testing.T) {
	th := throttle.New(50)
	th.Stop()
	th.Stop() // must not panic
}

func TestStop_Unlimited_NoPanic(t *testing.T) {
	th := throttle.New(0)
	th.Stop()
	th.Stop()
}

func TestNew_NegativeRPS_TreatedAsUnlimited(t *testing.T) {
	th := throttle.New(-10)
	defer th.Stop()

	ctx := context.Background()
	done := make(chan error, 1)
	go func() { done <- th.Wait(ctx) }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatal("Wait blocked for negative-RPS throttle")
	}
}
