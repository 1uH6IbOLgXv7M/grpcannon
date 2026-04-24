// Package loadgen orchestrates a timed load generation run, driving
// workers at a concurrency level supplied by a profile stage channel
// and collecting per-request metrics into a Recorder.
package loadgen

import (
	"context"
	"sync"
	"time"

	"github.com/example/grpcannon/internal/metrics"
	"github.com/example/grpcannon/internal/throttle"
)

// RequestFunc is the unit of work executed by each worker.
type RequestFunc func(ctx context.Context) error

// Config holds the parameters for a single load generation run.
type Config struct {
	// Stages is a channel of concurrency levels produced by a profile.
	Stages <-chan int
	// RPS is the maximum requests per second (0 = unlimited).
	RPS int
	// Fn is the gRPC call to execute.
	Fn RequestFunc
	// Recorder accumulates latency and error observations.
	Recorder *metrics.Recorder
}

// Run drives load according to cfg until Stages is closed or ctx is cancelled.
// It adjusts the worker pool size whenever a new stage value arrives.
func Run(ctx context.Context, cfg Config) error {
	th := throttle.New(cfg.RPS)
	defer th.Stop()

	var (
		mu      sync.Mutex
		cancel  context.CancelFunc
		wg      sync.WaitGroup
		current int
	)

	startWorkers := func(n int) {
		mu.Lock()
		defer mu.Unlock()
		if cancel != nil {
			cancel()
		}
		wg.Wait()
		var workerCtx context.Context
		workerCtx, cancel = context.WithCancel(ctx)
		current = n
		for i := 0; i < current; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					if err := th.Wait(workerCtx); err != nil {
						return
					}
					start := time.Now()
					err := cfg.Fn(workerCtx)
					cfg.Recorder.Observe(time.Since(start), err)
				}
			}()
		}
	}

	for {
		select {
		case <-ctx.Done():
			mu.Lock()
			if cancel != nil {
				cancel()
			}
			mu.Unlock()
			wg.Wait()
			return ctx.Err()
		case n, ok := <-cfg.Stages:
			if !ok {
				mu.Lock()
				if cancel != nil {
					cancel()
				}
				mu.Unlock()
				wg.Wait()
				return nil
			}
			if n != current {
				startWorkers(n)
			}
		}
	}
}
