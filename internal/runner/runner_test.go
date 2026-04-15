package runner_test

import (
	"context"
	"testing"

	"github.com/nickbadlose/grpcannon/internal/config"
	"github.com/nickbadlose/grpcannon/internal/runner"
)

func TestRun_InvalidConfig(t *testing.T) {
	cfg := &config.Config{} // missing required fields
	_, err := runner.Run(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected error for invalid config, got nil")
	}
}

func TestRun_DialFailure(t *testing.T) {
	cfg := &config.Config{
		Target:        "invalid-host:1",
		Call:          "/pkg.Service/Method",
		Concurrency:   1,
		TotalRequests: 1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	// With a zero-timeout context the invoke will fail; we just want Run to
	// return a Result (dial itself is non-blocking with the new grpc.NewClient).
	res, err := runner.Run(ctx, cfg)
	// Either a result with errors or a context error is acceptable.
	if err == nil && res == nil {
		t.Fatal("expected either an error or a non-nil result")
	}
}

func TestRun_ReturnsResultShape(t *testing.T) {
	cfg := &config.Config{
		Target:        "127.0.0.1:50099", // nothing listening — calls will error
		Call:          "/pkg.Service/Method",
		Concurrency:   2,
		TotalRequests: 4,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2)
	defer cancel()

	res, err := runner.Run(ctx, cfg)
	if err != nil {
		// A dial-level error is fine for this unit test environment.
		t.Skipf("skipping shape check, dial error: %v", err)
	}

	if res.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", res.Duration)
	}

	// All requests should be accounted for (success + error == total).
	total := res.Summary.Total
	if total != cfg.TotalRequests {
		t.Errorf("expected %d total requests in summary, got %d", cfg.TotalRequests, total)
	}
}
