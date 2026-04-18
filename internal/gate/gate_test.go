package gate_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/grpcannon/internal/gate"
)

func TestNew_Unbounded_NeverBlocks(t *testing.T) {
	g := gate.New(0)
	ctx := context.Background()
	for i := 0; i < 100; i++ {
		if err := g.Wait(ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}

func TestWait_BlocksAtLimit(t *testing.T) {
	g := gate.New(1)
	ctx := context.Background()

	if err := g.Wait(ctx); err != nil {
		t.Fatal(err)
	}

	ctx2, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	if err := g.Wait(ctx2); err == nil {
		t.Fatal("expected error when gate is full")
	}
}

func TestDone_ReleasesSlot(t *testing.T) {
	g := gate.New(1)
	ctx := context.Background()

	if err := g.Wait(ctx); err != nil {
		t.Fatal(err)
	}
	g.Done()

	if err := g.Wait(ctx); err != nil {
		t.Fatalf("expected slot after Done: %v", err)
	}
}

func TestClose_ReturnsErrClosed(t *testing.T) {
	g := gate.New(2)
	g.Close()

	err := g.Wait(context.Background())
	if err != gate.ErrClosed {
		t.Fatalf("expected ErrClosed, got %v", err)
	}
}

func TestConcurrent_NeverExceedsLimit(t *testing.T) {
	const limit = 4
	const goroutines = 20

	g := gate.New(limit)
	var (
		mu      sync.Mutex
		current int
		max     int
		wg      sync.WaitGroup
	)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := g.Wait(context.Background()); err != nil {
				return
			}
			defer g.Done()
			mu.Lock()
			current++
			if current > max {
				max = current
			}
			mu.Unlock()
			time.Sleep(2 * time.Millisecond)
			mu.Lock()
			current--
			mu.Unlock()
		}()
	}
	wg.Wait()

	if max > limit {
		t.Fatalf("concurrency %d exceeded limit %d", max, limit)
	}
}
