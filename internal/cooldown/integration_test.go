package cooldown_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/cooldown"
)

// TestConcurrent_AllowIsSafe verifies that concurrent calls to Allow never
// allow more activations than the interval permits.
func TestConcurrent_AllowIsSafe(t *testing.T) {
	const goroutines = 50
	c := cooldown.New(5 * time.Millisecond)

	var allowed atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if c.Allow() {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()

	// Only one goroutine should have been allowed in the first burst.
	if got := allowed.Load(); got != 1 {
		t.Fatalf("expected exactly 1 allowed, got %d", got)
	}
}
