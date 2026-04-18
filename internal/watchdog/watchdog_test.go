package watchdog_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/watchdog"
)

func TestRun_ContextCancelled_ReturnsNil(t *testing.T) {
	var c watchdog.Counter
	wd := watchdog.New(watchdog.Config{Threshold: 0.1, Window: 100 * time.Millisecond}, &c)
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(nil)
	if err := wd.Run(ctx, cancel); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRun_ThresholdBreached_ReturnsError(t *testing.T) {
	var c watchdog.Counter
	for i := 0; i < 20; i++ {
		c.RecordError()
	}
	wd := watchdog.New(watchdog.Config{
		Threshold:   0.05,
		Window:      50 * time.Millisecond,
		MinRequests: 5,
	}, &c)
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)
	err := wd.Run(ctx, cancel)
	if err != watchdog.ErrThresholdExceeded {
		t.Fatalf("expected ErrThresholdExceeded, got %v", err)
	}
}

func TestRun_BelowMinRequests_DoesNotCancel(t *testing.T) {
	var c watchdog.Counter
	c.RecordError() // 1 error, below MinRequests=10
	wd := watchdog.New(watchdog.Config{
		Threshold:   0.0,
		Window:      100 * time.Millisecond,
		MinRequests: 10,
	}, &c)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	causeCancel := func(err error) { cancel() }
	_ = causeCancel
	ctxCC, cc := context.WithCancelCause(ctx)
	defer cc(nil)
	err := wd.Run(ctxCC, cc)
	if err != nil {
		t.Fatalf("expected nil (timeout exit), got %v", err)
	}
}

func TestCounter_ErrorRate(t *testing.T) {
	var c watchdog.Counter
	if c.ErrorRate() != 0 {
		t.Fatal("expected 0 for empty counter")
	}
	c.RecordSuccess()
	c.RecordSuccess()
	c.RecordError()
	got := c.ErrorRate()
	want := 1.0 / 3.0
	if got < want-0.001 || got > want+0.001 {
		t.Fatalf("expected ~%.3f, got %.3f", want, got)
	}
}

func TestCounter_Total(t *testing.T) {
	var c watchdog.Counter
	for i := 0; i < 5; i++ {
		c.RecordSuccess()
	}
	c.RecordError()
	if c.Total() != 6 {
		t.Fatalf("expected 6, got %d", c.Total())
	}
}
