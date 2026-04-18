package watchdog_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/watchdog"
)

func TestConcurrent_RecordAndWatch(t *testing.T) {
	var c watchdog.Counter
	wd := watchdog.New(watchdog.Config{
		Threshold:   0.5,
		Window:      50 * time.Millisecond,
		MinRequests: 10,
	}, &c)

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	var wg sync.WaitGroup
	wg.Add(1)
	var runErr error
	go func() {
		defer wg.Done()
		runErr = wd.Run(ctx, cancel)
	}()

	// Drive error rate above threshold.
	for i := 0; i < 30; i++ {
		c.RecordError()
	}

	wg.Wait()
	if !errors.Is(runErr, watchdog.ErrThresholdExceeded) {
		t.Fatalf("expected ErrThresholdExceeded, got %v", runErr)
	}
	if cause := context.Cause(ctx); !errors.Is(cause, watchdog.ErrThresholdExceeded) {
		t.Fatalf("expected context cause ErrThresholdExceeded, got %v", cause)
	}
}
