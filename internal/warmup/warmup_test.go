package warmup_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/warmup"
)

type mockDoer struct {
	calls atomic.Int64
	errAfter int64
}

func (m *mockDoer) Do(_ context.Context, _ string, _ []byte) error {
	n := m.calls.Add(1)
	if m.errAfter > 0 && n > m.errAfter {
		return errors.New("mock error")
	}
	return nil
}

func TestRun_ZeroRequests_NoCallsMade(t *testing.T) {
	d := &mockDoer{}
	cfg := warmup.Config{Requests: 0, Concurrency: 1, Timeout: time.Second}
	errCount, err := warmup.Run(context.Background(), cfg, "/svc/Method", nil, d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if errCount != 0 {
		t.Fatalf("expected 0 errors, got %d", errCount)
	}
	if d.calls.Load() != 0 {
		t.Fatalf("expected 0 calls, got %d", d.calls.Load())
	}
}

func TestRun_AllSucceed_ZeroErrors(t *testing.T) {
	d := &mockDoer{}
	cfg := warmup.Config{Requests: 8, Concurrency: 4, Timeout: time.Second}
	errCount, err := warmup.Run(context.Background(), cfg, "/svc/Method", nil, d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if errCount != 0 {
		t.Fatalf("expected 0 errors, got %d", errCount)
	}
	if d.calls.Load() != 8 {
		t.Fatalf("expected 8 calls, got %d", d.calls.Load())
	}
}

func TestRun_SomeErrors_CountedCorrectly(t *testing.T) {
	d := &mockDoer{errAfter: 5}
	cfg := warmup.Config{Requests: 10, Concurrency: 2, Timeout: time.Second}
	errCount, err := warmup.Run(context.Background(), cfg, "/svc/Method", nil, d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if errCount != 5 {
		t.Fatalf("expected 5 errors, got %d", errCount)
	}
}

func TestRun_ContextCancelled_ReturnsError(t *testing.T) {
	blocking := &blockingDoer{}
	cfg := warmup.Config{Requests: 4, Concurrency: 2, Timeout: time.Second}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	_, err := warmup.Run(ctx, cfg, "/svc/Method", nil, blocking)
	if err == nil {
		t.Fatal("expected error due to cancellation")
	}
}

type blockingDoer struct{}

func (b *blockingDoer) Do(ctx context.Context, _ string, _ []byte) error {
	<-ctx.Done()
	return ctx.Err()
}

func TestDefault_ReturnsNonZeroValues(t *testing.T) {
	cfg := warmup.Default()
	if cfg.Requests <= 0 {
		t.Error("expected positive Requests")
	}
	if cfg.Concurrency <= 0 {
		t.Error("expected positive Concurrency")
	}
	if cfg.Timeout <= 0 {
		t.Error("expected positive Timeout")
	}
}
