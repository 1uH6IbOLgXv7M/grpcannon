package pause

import (
	"sync"
	"testing"
	"time"
)

func TestNew_IsResumedByDefault(t *testing.T) {
	c := New()
	if c.IsPaused() {
		t.Fatal("expected controller to be resumed by default")
	}
}

func TestPause_SetsIsPaused(t *testing.T) {
	c := New()
	c.Pause()
	if !c.IsPaused() {
		t.Fatal("expected controller to be paused")
	}
}

func TestResume_ClearsIsPaused(t *testing.T) {
	c := New()
	c.Pause()
	c.Resume()
	if c.IsPaused() {
		t.Fatal("expected controller to be resumed")
	}
}

func TestWait_ReturnsImmediatelyWhenResumed(t *testing.T) {
	c := New()
	done := make(chan struct{})
	go func() {
		c.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Wait did not return promptly when not paused")
	}
}

func TestWait_BlocksWhilePaused(t *testing.T) {
	c := New()
	c.Pause()
	blocked := make(chan struct{})
	go func() {
		close(blocked)
		c.Wait()
	}()
	<-blocked
	time.Sleep(20 * time.Millisecond)
	if !c.IsPaused() {
		t.Fatal("controller should still be paused")
	}
	c.Resume()
}

func TestResume_UnblocksMultipleWaiters(t *testing.T) {
	c := New()
	c.Pause()
	const n = 10
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			c.Wait()
			wg.Done()
		}()
	}
	time.Sleep(20 * time.Millisecond)
	c.Resume()
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("not all waiters were unblocked after Resume")
	}
}

func TestPause_CalledTwice_NoPanic(t *testing.T) {
	c := New()
	c.Pause()
	c.Pause()
}

func TestResume_CalledTwice_NoPanic(t *testing.T) {
	c := New()
	c.Resume()
	c.Resume()
}
