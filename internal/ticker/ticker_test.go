package ticker_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/ticker"
)

func TestNew_DefaultInterval_WhenZero(t *testing.T) {
	tk := ticker.New(0)
	if tk.Interval() != time.Second {
		t.Fatalf("expected 1s default, got %v", tk.Interval())
	}
}

func TestNew_DefaultInterval_WhenNegative(t *testing.T) {
	tk := ticker.New(-5 * time.Millisecond)
	if tk.Interval() != time.Second {
		t.Fatalf("expected 1s default, got %v", tk.Interval())
	}
}

func TestNew_CustomInterval_Preserved(t *testing.T) {
	tk := ticker.New(50 * time.Millisecond)
	if tk.Interval() != 50*time.Millisecond {
		t.Fatalf("unexpected interval %v", tk.Interval())
	}
}

func TestRun_DeliversAtLeastOneTick(t *testing.T) {
	tk := ticker.New(20 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go tk.Run(ctx)

	select {
	case _, ok := <-tk.C():
		if !ok {
			t.Fatal("channel closed before tick")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for tick")
	}
}

func TestRun_ChannelClosedAfterCancel(t *testing.T) {
	tk := ticker.New(10 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		tk.Run(ctx)
		close(done)
	}()

	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not exit after cancel")
	}

	// Channel must be closed.
	for range tk.C() {
	}
}

func TestRun_MultipleTicksDelivered(t *testing.T) {
	tk := ticker.New(15 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go tk.Run(ctx)

	count := 0
	timeout := time.After(300 * time.Millisecond)
loop:
	for {
		select {
		case _, ok := <-tk.C():
			if !ok {
				break loop
			}
			count++
			if count >= 3 {
				break loop
			}
		case <-timeout:
			break loop
		}
	}

	if count < 3 {
		t.Fatalf("expected at least 3 ticks, got %d", count)
	}
}
