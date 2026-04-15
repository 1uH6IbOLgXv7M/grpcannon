package throttle_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/throttle"
)

// TestRateApproximation verifies that over a short window the throttle
// delivers roughly the requested number of tokens.
func TestRateApproximation(t *testing.T) {
	const rps = 200
	const window = 300 * time.Millisecond

	th := throttle.New(rps)
	defer th.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), window)
	defer cancel()

	var count int64
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if err := th.Wait(ctx); err != nil {
				return
			}
			atomic.AddInt64(&count, 1)
		}
	}()

	<-done

	got := atomic.LoadInt64(&count)
	// Allow ±40 % tolerance for CI timing jitter.
	expected := int64(float64(rps) * window.Seconds())
	lo := int64(float64(expected) * 0.60)
	hi := int64(float64(expected) * 1.40)
	if got < lo || got > hi {
		t.Errorf("token count %d outside [%d, %d] for %d rps over %s",
			got, lo, hi, rps, window)
	}
}
