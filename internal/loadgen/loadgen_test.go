package loadgen_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/loadgen"
	"github.com/example/grpcannon/internal/metrics"
)

func stagesOf(vals ...int) <-chan int {
	ch := make(chan int, len(vals))
	for _, v := range vals {
		ch <- v
	}
	close(ch)
	return ch
}

func TestRun_ClosedStages_ReturnsNil(t *testing.T) {
	rec := metrics.NewRecorder()
	err := loadgen.Run(context.Background(), loadgen.Config{
		Stages:   stagesOf(),
		Recorder: rec,
		Fn: func(_ context.Context) error {
			return nil
		},
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRun_ContextCancelled_ReturnsError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan int) // never closed
	rec := metrics.NewRecorder()

	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	err := loadgen.Run(ctx, loadgen.Config{
		Stages:   ch,
		Recorder: rec,
		Fn: func(_ context.Context) error { return nil },
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestRun_RecordsRequests(t *testing.T) {
	var count atomic.Int64
	rec := metrics.NewRecorder()

	ch := make(chan int, 1)
	ch <- 4

	go func() {
		time.Sleep(50 * time.Millisecond)
		close(ch)
	}()

	_ = loadgen.Run(context.Background(), loadgen.Config{
		Stages:   ch,
		RPS:      0,
		Recorder: rec,
		Fn: func(ctx context.Context) error {
			count.Add(1)
			time.Sleep(2 * time.Millisecond)
			return nil
		},
	})

	snap := rec.Snapshot()
	if snap.Total == 0 {
		t.Fatal("expected at least one recorded request")
	}
	if snap.Total != count.Load() {
		t.Fatalf("recorder total %d != fn calls %d", snap.Total, count.Load())
	}
}

func TestRun_RecordsErrors(t *testing.T) {
	sentinel := errors.New("boom")
	rec := metrics.NewRecorder()

	ch := make(chan int, 1)
	ch <- 2

	go func() {
		time.Sleep(40 * time.Millisecond)
		close(ch)
	}()

	_ = loadgen.Run(context.Background(), loadgen.Config{
		Stages:   ch,
		Recorder: rec,
		Fn: func(_ context.Context) error {
			time.Sleep(3 * time.Millisecond)
			return sentinel
		},
	})

	snap := rec.Snapshot()
	if snap.Errors == 0 {
		t.Fatal("expected errors to be recorded")
	}
	if snap.Errors != snap.Total {
		t.Fatalf("all requests should be errors: total=%d errors=%d", snap.Total, snap.Errors)
	}
}
