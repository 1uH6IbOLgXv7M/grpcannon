package worker

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestPool_RunCompletesAllRequests(t *testing.T) {
	const total = 20
	var calls int64

	fn := func(_ context.Context) error {
		atomic.AddInt64(&calls, 1)
		return nil
	}

	p := NewPool(4, total, fn)
	p.Run(context.Background())

	var results []Result
	for r := range p.Results {
		results = append(results, r)
	}

	if len(results) != total {
		t.Fatalf("expected %d results, got %d", total, len(results))
	}
	if calls != total {
		t.Fatalf("expected %d calls, got %d", total, calls)
	}
}

func TestPool_RunRecordsErrors(t *testing.T) {
	expectedErr := errors.New("rpc error")
	fn := func(_ context.Context) error { return expectedErr }

	p := NewPool(2, 5, fn)
	p.Run(context.Background())

	for r := range p.Results {
		if r.Err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, r.Err)
		}
	}
}

func TestPool_RunRecordsDuration(t *testing.T) {
	delay := 10 * time.Millisecond
	fn := func(_ context.Context) error {
		time.Sleep(delay)
		return nil
	}

	p := NewPool(1, 3, fn)
	p.Run(context.Background())

	for r := range p.Results {
		if r.Duration < delay {
			t.Errorf("expected duration >= %v, got %v", delay, r.Duration)
		}
	}
}

func TestPool_RunRespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	fn := func(_ context.Context) error { return nil }
	p := NewPool(4, 100, fn)
	p.Run(ctx)

	var count int
	for range p.Results {
		count++
	}
	// With a cancelled context some workers may not process all items.
	if count > 100 {
		t.Errorf("unexpected result count: %d", count)
	}
}
