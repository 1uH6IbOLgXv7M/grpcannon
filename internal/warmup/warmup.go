// Package warmup provides a pre-load warm-up phase that sends a small
// number of requests before the main test begins, allowing JIT compilation
// and connection establishment to stabilise.
package warmup

import (
	"context"
	"fmt"
	"time"
)

// Doer is the interface satisfied by invoker.Invoker.
type Doer interface {
	Do(ctx context.Context, method string, payload []byte) error
}

// Config holds warm-up parameters.
type Config struct {
	// Requests is the number of warm-up requests to send.
	Requests int
	// Concurrency is the number of parallel goroutines used during warm-up.
	Concurrency int
	// Timeout is the per-request deadline.
	Timeout time.Duration
}

// Default returns a Config with sensible defaults.
func Default() Config {
	return Config{
		Requests:    10,
		Concurrency: 2,
		Timeout:     5 * time.Second,
	}
}

// Run executes the warm-up phase, returning the number of errors encountered.
// It does not fail on individual request errors; callers may inspect the count.
func Run(ctx context.Context, cfg Config, method string, payload []byte, doer Doer) (int, error) {
	if cfg.Requests <= 0 {
		return 0, nil
	}
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 1
	}
	// Clamp concurrency to the number of requests to avoid idle goroutines.
	if cfg.Concurrency > cfg.Requests {
		cfg.Concurrency = cfg.Requests
	}

	work := make(chan struct{}, cfg.Requests)
	for i := 0; i < cfg.Requests; i++ {
		work <- struct{}{}
	}
	close(work)

	type result struct{ err error }
	results := make(chan result, cfg.Requests)

	for w := 0; w < cfg.Concurrency; w++ {
		go func() {
			for range work {
				reqCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
				err := doer.Do(reqCtx, method, payload)
				cancel()
				results <- result{err: err}
			}
		}()
	}

	errCount := 0
	for i := 0; i < cfg.Requests; i++ {
		select {
		case r := <-results:
			if r.err != nil {
				errCount++
			}
			case <-ctx.Done():
			return errCount, fmt.Errorf("warmup cancelled: %w", ctx.Err())
		}
	}
	return errCount, nil
}
