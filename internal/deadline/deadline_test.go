package deadline_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/deadline"
)

func TestNew_ZeroTimeout_NoDeadlineSet(t *testing.T) {
	e := deadline.New(0)
	ctx, cancel := e.Wrap(context.Background())
	defer cancel()
	if _, ok := ctx.Deadline(); ok {
		t.Fatal("expected no deadline for zero timeout")
	}
}

func TestNew_NegativeTimeout_NoDeadlineSet(t *testing.T) {
	e := deadline.New(-1 * time.Second)
	ctx, cancel := e.Wrap(context.Background())
	defer cancel()
	if _, ok := ctx.Deadline(); ok {
		t.Fatal("expected no deadline for negative timeout")
	}
}

func TestNew_PositiveTimeout_DeadlineSet(t *testing.T) {
	e := deadline.New(5 * time.Second)
	ctx, cancel := e.Wrap(context.Background())
	defer cancel()
	if _, ok := ctx.Deadline(); !ok {
		t.Fatal("expected deadline to be set")
	}
}

func TestWrap_ContextExpiresAfterTimeout(t *testing.T) {
	e := deadline.New(20 * time.Millisecond)
	ctx, cancel := e.Wrap(context.Background())
	defer cancel()
	select {
	case <-ctx.Done():
		// expected
	case <-time.After(500 * time.Millisecond):
		t.Fatal("context did not expire in time")
	}
}

func TestWrap_CancelFuncCancelsContext(t *testing.T) {
	e := deadline.New(0)
	ctx, cancel := e.Wrap(context.Background())
	cancel()
	select {
	case <-ctx.Done():
	default:
		t.Fatal("expected context to be cancelled")
	}
}

func TestIsExceeded_DeadlineExceeded(t *testing.T) {
	if !deadline.IsExceeded(context.DeadlineExceeded) {
		t.Fatal("expected true for context.DeadlineExceeded")
	}
}

func TestIsExceeded_ErrExceeded(t *testing.T) {
	if !deadline.IsExceeded(deadline.ErrExceeded) {
		t.Fatal("expected true for deadline.ErrExceeded")
	}
}

func TestIsExceeded_OtherError(t *testing.T) {
	if deadline.IsExceeded(context.Canceled) {
		t.Fatal("expected false for context.Canceled")
	}
}
