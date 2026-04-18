package cooldown_test

import (
	"testing"
	"time"

	"github.com/example/grpcannon/internal/cooldown"
)

func TestNew_DefaultInterval_WhenZero(t *testing.T) {
	c := cooldown.New(0)
	if c == nil {
		t.Fatal("expected non-nil cooldown")
	}
}

func TestAllow_FirstCall_ReturnsTrue(t *testing.T) {
	c := cooldown.New(time.Second)
	if !c.Allow() {
		t.Fatal("expected first Allow to return true")
	}
}

func TestAllow_ImmediateSecondCall_ReturnsFalse(t *testing.T) {
	c := cooldown.New(time.Second)
	c.Allow()
	if c.Allow() {
		t.Fatal("expected second immediate Allow to return false")
	}
}

func TestAllow_AfterInterval_ReturnsTrue(t *testing.T) {
	c := cooldown.New(10 * time.Millisecond)
	c.Allow()
	time.Sleep(20 * time.Millisecond)
	if !c.Allow() {
		t.Fatal("expected Allow to return true after interval")
	}
}

func TestRemaining_BeforeFirstAllow_ReturnsZero(t *testing.T) {
	c := cooldown.New(time.Second)
	if r := c.Remaining(); r != 0 {
		t.Fatalf("expected 0, got %v", r)
	}
}

func TestRemaining_AfterAllow_IsPositive(t *testing.T) {
	c := cooldown.New(time.Second)
	c.Allow()
	if r := c.Remaining(); r <= 0 {
		t.Fatalf("expected positive remaining, got %v", r)
	}
}

func TestRemaining_AfterInterval_ReturnsZero(t *testing.T) {
	c := cooldown.New(10 * time.Millisecond)
	c.Allow()
	time.Sleep(20 * time.Millisecond)
	if r := c.Remaining(); r != 0 {
		t.Fatalf("expected 0 after interval, got %v", r)
	}
}

func TestReset_AllowsImmediateReactivation(t *testing.T) {
	c := cooldown.New(time.Hour)
	c.Allow()
	c.Reset()
	if !c.Allow() {
		t.Fatal("expected Allow to return true after Reset")
	}
}
