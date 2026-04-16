package signal_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/grpcannon/grpcannon/internal/signal"
)

func TestNotifyContext_CancelledByParent(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	ctx, stop := signal.NotifyContext(parent)
	defer stop()

	cancel()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled after parent cancel")
	}
}

func TestNotifyContext_StopReleasesResources(t *testing.T) {
	ctx, stop := signal.NotifyContext(context.Background())
	stop()

	// After stop the context should not be cancelled on its own.
	select {
	case <-ctx.Done():
		t.Fatal("context should not be done after stop without signal")
	default:
		// expected
	}
}

func TestTrap_CancelsOnSignal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		signal.Trap(cancel)
		close(done)
	}()

	// Give the goroutine time to register.
	time.Sleep(20 * time.Millisecond)

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("find process: %v", err)
	}
	if err := p.Signal(os.Interrupt); err != nil {
		t.Fatalf("send signal: %v", err)
	}

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(2 * time.Second):
		t.Fatal("context was not cancelled after SIGINT")
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Trap goroutine did not exit")
	}
}
