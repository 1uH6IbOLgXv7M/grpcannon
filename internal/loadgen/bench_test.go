package loadgen_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/loadgen"
	"github.com/example/grpcannon/internal/metrics"
)

// BenchmarkRun_HighConcurrency measures the per-request overhead of the
// loadgen loop itself using a no-op RequestFunc.
func BenchmarkRun_HighConcurrency(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rec := metrics.NewRecorder()
		ch := make(chan int, 1)
		ch <- 8

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		_ = loadgen.Run(ctx, loadgen.Config{
			Stages:   ch,
			Recorder: rec,
			Fn:       func(_ context.Context) error { return nil },
		})
		cancel()
	}
}
