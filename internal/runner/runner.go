// Package runner orchestrates a load test run by wiring together
// the worker pool, metrics recorder, and configuration.
package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/nickbadlose/grpcannon/internal/config"
	"github.com/nickbadlose/grpcannon/internal/metrics"
	"github.com/nickbadlose/grpcannon/internal/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Result holds the outcome of a completed load test.
type Result struct {
	Summary  metrics.Summary
	Duration time.Duration
}

// Run executes a load test according to cfg and returns aggregated results.
func Run(ctx context.Context, cfg *config.Config) (*Result, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	conn, err := grpc.NewClient(
		cfg.Target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", cfg.Target, err)
	}
	defer conn.Close()

	rec := metrics.NewRecorder()

	callFn := func(ctx context.Context) error {
		return conn.Invoke(ctx, cfg.Call, cfg.RequestData, nil)
	}

	pool := worker.NewPool(cfg.Concurrency, cfg.TotalRequests, callFn, rec)

	start := time.Now()
	pool.Run(ctx)
	elapsed := time.Since(start)

	snap := rec.Snapshot()

	return &Result{
		Summary:  snap,
		Duration: elapsed,
	}, nil
}
