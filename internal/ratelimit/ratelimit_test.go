package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/ratelimit"
)

func TestNew_Unlimited_DoesNotBlock(t *testing.T) {
	l := ratelimit.New(ratelimit.Unlimited)
	defer l.Stop()

	ctx := context.Background()
	done := make(chan error, 1)
	go func() { done <- l.Wait(ctx) }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatal("Wait blocked unexpectedly for unlimited limiter")
	}
}

func TestNew_RateLimited_AllowsTokens(t *testing.T) {
	const rps = 100
	l := ratelimit.New(rps)
	defer l.Stop()

	ctx := context.Background()
	start := time.Now()
	for i := 0; i < 5; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("unexpected error on token %d: %v", i, err)
		}
	}
	elapsed := time.Since(start)

	// 5 tokens at 100 rps → ~40 ms minimum spacing; allow generous upper bound.
	if elapsed < 30*time.Millisecond {
		t.Errorf("tokens arrived too fast: %v", elapsed)
	}
}

func TestWait_ContextCancelled_ReturnsError(t *testing.T) {
	// 1 rps so the ticker fires every second — context should cancel first.
	l := ratelimit.New(1)
	defer l.Stop()

	// Drain the first tick that fires immediately.
	ctx := context.Background()
	_ = l.Wait(ctx)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := l.Wait(ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestStop_CalledTwice_NoPanic(t *testing.T) {
	l := ratelimit.New(50)
	l.Stop()
	l.Stop() // should not panic
}
