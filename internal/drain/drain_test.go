package drain_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/drain"
)

func TestAcquire_BeforeClose_ReturnsTrue(t *testing.T) {
	d := drain.New()
	if !d.Acquire() {
		t.Fatal("expected Acquire to return true before Drain")
	}
	d.Release()
}

func TestAcquire_AfterDrain_ReturnsFalse(t *testing.T) {
	d := drain.New()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := d.Drain(ctx); err != nil {
		t.Fatalf("unexpected drain error: %v", err)
	}
	if d.Acquire() {
		t.Fatal("expected Acquire to return false after Drain")
	}
}

func TestDrain_WaitsForInFlight(t *testing.T) {
	d := drain.New()
	if !d.Acquire() {
		t.Fatal("acquire failed")
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		d.Release()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := d.Drain(ctx); err != nil {
		t.Fatalf("drain returned error: %v", err)
	}
}

func TestDrain_ContextCancelled_ReturnsError(t *testing.T) {
	d := drain.New()
	if !d.Acquire() {
		t.Fatal("acquire failed")
	}
	defer d.Release()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if err := d.Drain(ctx); err == nil {
		t.Fatal("expected error when context expires")
	}
}

func TestDrain_ConcurrentAcquireRelease(t *testing.T) {
	d := drain.New()
	const n = 50
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if d.Acquire() {
				time.Sleep(5 * time.Millisecond)
				d.Release()
			}
		}()
	}
	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := d.Drain(ctx); err != nil {
		t.Fatalf("unexpected drain error: %v", err)
	}
}

func TestDrainTimeout_Convenience(t *testing.T) {
	d := drain.New()
	if err := d.DrainTimeout(time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
