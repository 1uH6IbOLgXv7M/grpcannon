package debounce_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/nickpoorman/grpcannon/internal/debounce"
)

func TestNew_DefaultInterval_WhenZero(t *testing.T) {
	var calls int32
	d := debounce.New(0, func() { atomic.AddInt32(&calls, 1) })
	d.Call()
	time.Sleep(200 * time.Millisecond)
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("expected 1 call, got %d", got)
	}
}

func TestCall_FiresAfterInterval(t *testing.T) {
	var calls int32
	d := debounce.New(50*time.Millisecond, func() { atomic.AddInt32(&calls, 1) })
	d.Call()
	time.Sleep(120 * time.Millisecond)
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestCall_ResetsTimer(t *testing.T) {
	var calls int32
	d := debounce.New(80*time.Millisecond, func() { atomic.AddInt32(&calls, 1) })
	d.Call()
	time.Sleep(40 * time.Millisecond)
	d.Call() // reset
	time.Sleep(40 * time.Millisecond)
	// should not have fired yet
	if got := atomic.LoadInt32(&calls); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
	time.Sleep(80 * time.Millisecond)
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestFlush_FiresImmediately(t *testing.T) {
	var calls int32
	d := debounce.New(500*time.Millisecond, func() { atomic.AddInt32(&calls, 1) })
	d.Call()
	flushed := d.Flush()
	if !flushed {
		t.Fatal("expected Flush to return true")
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestFlush_NoPending_ReturnsFalse(t *testing.T) {
	d := debounce.New(50*time.Millisecond, func() {})
	if d.Flush() {
		t.Fatal("expected false when nothing pending")
	}
}

func TestStop_CancelsPending(t *testing.T) {
	var calls int32
	d := debounce.New(50*time.Millisecond, func() { atomic.AddInt32(&calls, 1) })
	d.Call()
	d.Stop()
	time.Sleep(100 * time.Millisecond)
	if got := atomic.LoadInt32(&calls); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestStop_CalledTwice_NoPanic(t *testing.T) {
	d := debounce.New(50*time.Millisecond, func() {})
	d.Call()
	d.Stop()
	d.Stop()
}
