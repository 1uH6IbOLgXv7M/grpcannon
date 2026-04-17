package fence

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNew_IsClosedByDefault(t *testing.T) {
	f := New()
	if f.IsOpen() {
		t.Fatal("expected fence to be closed after New")
	}
}

func TestOpen_SetsIsOpen(t *testing.T) {
	f := New()
	f.Open()
	if !f.IsOpen() {
		t.Fatal("expected fence to be open after Open")
	}
}

func TestOpen_CalledTwice_NoPanic(t *testing.T) {
	f := New()
	f.Open()
	f.Open() // must not panic
}

func TestWait_ReturnsImmediatelyWhenAlreadyOpen(t *testing.T) {
	f := New()
	f.Open()
	ctx := context.Background()
	if err := f.Wait(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWait_BlocksUntilOpen(t *testing.T) {
	f := New()
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := f.Wait(ctx); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	time.Sleep(20 * time.Millisecond)
	f.Open()
	wg.Wait()
}

func TestWait_ContextCancelled_ReturnsError(t *testing.T) {
	f := New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := f.Wait(ctx); err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestWait_ConcurrentWaiters_AllUnblocked(t *testing.T) {
	f := New()
	ctx := context.Background()
	const n = 20
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			_ = f.Wait(ctx)
		}()
	}
	time.Sleep(10 * time.Millisecond)
	f.Open()
	wg.Wait()
}
