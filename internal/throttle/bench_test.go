package throttle_test

import (
	"context"
	"testing"

	"github.com/example/grpcannon/internal/throttle"
)

// BenchmarkWait_Unlimited measures overhead of Wait when no rate limit is set.
func BenchmarkWait_Unlimited(b *testing.B) {
	th := throttle.New(0)
	defer th.Stop()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = th.Wait(ctx)
	}
}

// BenchmarkWait_HighRPS measures throughput at a very high token rate.
func BenchmarkWait_HighRPS(b *testing.B) {
	th := throttle.New(1_000_000)
	defer th.Stop()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = th.Wait(ctx)
	}
}
