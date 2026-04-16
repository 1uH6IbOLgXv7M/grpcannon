package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yourusername/grpcannon/internal/retry"
)

func TestDefault_SingleAttempt(t *testing.T) {
	p := retry.Default()
	if p.MaxAttempts != 1 {
		t.Fatalf("expected MaxAttempts=1, got %d", p.MaxAttempts)
	}
}

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.Default(), func(_ context.Context) error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesUpToMaxAttempts(t *testing.T) {
	p := retry.Policy{MaxAttempts: 3, Delay: 0}
	calls := 0
	sentinel := errors.New("transient")
	err := retry.Do(context.Background(), p, func(_ context.Context) error {
		calls++
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_StopsOnSuccess(t *testing.T) {
	p := retry.Policy{MaxAttempts: 5, Delay: 0}
	calls := 0
	err := retry.Do(context.Background(), p, func(_ context.Context) error {
		calls++
		if calls < 3 {
			return errors.New("not yet")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_RetryOn_SkipsNonMatchingCode(t *testing.T) {
	p := retry.Policy{
		MaxAttempts: 3,
		RetryOn:     []codes.Code{codes.Unavailable},
	}
	calls := 0
	permErr := status.Error(codes.InvalidArgument, "bad")
	err := retry.Do(context.Background(), p, func(_ context.Context) error {
		calls++
		return permErr
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if calls != 1 {
		t.Fatalf("expected 1 call (no retry), got %d", calls)
	}
}

func TestDo_RetryOn_RetriesMatchingCode(t *testing.T) {
	p := retry.Policy{
		MaxAttempts: 3,
		RetryOn:     []codes.Code{codes.Unavailable},
	}
	calls := 0
	err := retry.Do(context.Background(), p, func(_ context.Context) error {
		calls++
		return status.Error(codes.Unavailable, "down")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ContextCancelled_StopsEarly(t *testing.T) {
	p := retry.Policy{MaxAttempts: 10, Delay: 50 * time.Millisecond}
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err := retry.Do(ctx, p, func(_ context.Context) error {
		calls++
		return errors.New("fail")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls >= 10 {
		t.Fatal("expected early stop due to cancellation")
	}
}
