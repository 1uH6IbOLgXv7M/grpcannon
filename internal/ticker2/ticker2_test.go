package ticker2_test

import (
	"testing"
	"time"

	"github.com/example/grpcannon/internal/ticker2"
)

func TestNew_DefaultInterval_WhenZero(t *testing.T) {
	tk := ticker2.New(0)
	defer tk.Stop()
	// channel must be non-nil
	if tk.C == nil {
		t.Fatal("expected non-nil channel")
	}
}

func TestNew_DefaultInterval_WhenNegative(t *testing.T) {
	tk := ticker2.New(-5 * time.Second)
	defer tk.Stop()
	if tk.C == nil {
		t.Fatal("expected non-nil channel")
	}
}

func TestNew_CustomInterval_Fires(t *testing.T) {
	tk := ticker2.New(20 * time.Millisecond)
	defer tk.Stop()

	select {
	case <-tk.C:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("ticker did not fire within 500ms")
	}
}

func TestStop_SilencesChannel(t *testing.T) {
	tk := ticker2.New(10 * time.Millisecond)
	// drain one tick so the ticker has started
	<-tk.C
	tk.Stop()

	// After Stop, no more ticks should arrive quickly.
	time.Sleep(50 * time.Millisecond)
	select {
	case <-tk.C:
		t.Fatal("received tick after Stop")
	default:
	}
}

func TestReset_ChangesInterval(t *testing.T) {
	tk := ticker2.New(500 * time.Millisecond)
	defer tk.Stop()

	// Reset to a much shorter interval.
	tk.Reset(20 * time.Millisecond)

	select {
	case <-tk.C:
		// received tick after reset – good
	case <-time.After(300 * time.Millisecond):
		t.Fatal("ticker did not fire after Reset")
	}
}

func TestReset_ZeroInterval_Defaults(t *testing.T) {
	tk := ticker2.New(10 * time.Millisecond)
	defer tk.Stop()
	// Should not panic when reset with zero.
	tk.Reset(0)
}
