package limiter_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/your-org/grpcannon/internal/limiter"
)

func TestNew_Unbounded_AcquireNeverBlocks(t *testing.T) {
	l := limiter.New(0)
	for i := 0; i < 100; i++ {
		if err := l.Acquire(context.Background()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if got := l.Available(); got != -1 {
		t.Fatalf("expected -1, got %d", got)
	}
}

func TestAcquire_BlocksAtLimit(t *testing.T) {
	l := limiter.New(2)
	ctx := context.Background()
	_ = l.Acquire(ctx)
	_ = l.Acquire(ctx)

	ctx2, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if err := l.Acquire(ctx2); err == nil {
		t.Fatal("expected error when limit reached")
	}
}

func TestRelease_FreesSlot(t *testing.T) {
	l := limiter.New(1)
	_ = l.Acquire(context.Background())
	l.Release()
	if err := l.Acquire(context.Background()); err != nil {
		t.Fatalf("expected slot after release: %v", err)
	}
}

func TestAvailable_ReflectsUsage(t *testing.T) {
	l := limiter.New(3)
	if l.Available() != 3 {
		t.Fatalf("expected 3 available")
	}
	_ = l.Acquire(context.Background())
	if l.Available() != 2 {
		t.Fatalf("expected 2 available")
	}
}

func TestConcurrent_NeverExceedsLimit(t *testing.T) {
	const limit = 5
	l := limiter.New(limit)
	var mu sync.Mutex
	var peak, current int
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = l.Acquire(context.Background())
			mu.Lock()
			current++
			if current > peak {
				peak = current
			}
			mu.Unlock()
			time.Sleep(5 * time.Millisecond)
			mu.Lock()
			current--
			mu.Unlock()
			l.Release()
		}()
	}
	wg.Wait()
	if peak > limit {
		t.Fatalf("peak concurrency %d exceeded limit %d", peak, limit)
	}
}
