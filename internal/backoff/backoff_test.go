package backoff_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/grpcannon/internal/backoff"
)

func TestNone_NeverRetries(t *testing.T) {
	s := backoff.None{}
	_, ok := s.Next(0)
	if ok {
		t.Fatal("expected None to return ok=false")
	}
}

func TestConstant_ReturnsDelayWithinMaxRetries(t *testing.T) {
	s := backoff.Constant{Delay: 10 * time.Millisecond, MaxRetries: 3}
	for i := 0; i < 3; i++ {
		d, ok := s.Next(i)
		if !ok {
			t.Fatalf("attempt %d: expected ok=true", i)
		}
		if d != 10*time.Millisecond {
			t.Fatalf("attempt %d: expected 10ms, got %v", i, d)
		}
	}
}

func TestConstant_StopsAfterMaxRetries(t *testing.T) {
	s := backoff.Constant{Delay: 5 * time.Millisecond, MaxRetries: 2}
	_, ok := s.Next(2)
	if ok {
		t.Fatal("expected ok=false after MaxRetries")
	}
}

func TestExponential_GrowsWithAttempt(t *testing.T) {
	s := backoff.Exponential{
		BaseDelay:  10 * time.Millisecond,
		MaxDelay:   500 * time.Millisecond,
		MaxRetries: 5,
	}
	prev := time.Duration(0)
	for i := 0; i < 5; i++ {
		d, ok := s.Next(i)
		if !ok {
			t.Fatalf("attempt %d: expected ok=true", i)
		}
		if d < prev {
			t.Fatalf("attempt %d: delay %v decreased from %v", i, d, prev)
		}
		prev = d
	}
}

func TestExponential_CapsAtMaxDelay(t *testing.T) {
	s := backoff.Exponential{
		BaseDelay:  100 * time.Millisecond,
		MaxDelay:   200 * time.Millisecond,
		MaxRetries: 10,
	}
	for i := 3; i < 10; i++ {
		d, _ := s.Next(i)
		if d > 200*time.Millisecond {
			t.Fatalf("attempt %d: delay %v exceeds MaxDelay", i, d)
		}
	}
}

func TestExponential_StopsAfterMaxRetries(t *testing.T) {
	s := backoff.Exponential{BaseDelay: 1 * time.Millisecond, MaxDelay: 1 * time.Second, MaxRetries: 3}
	_, ok := s.Next(3)
	if ok {
		t.Fatal("expected ok=false after MaxRetries")
	}
}

func TestWait_ContextCancelled_ReturnsError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	s := backoff.Constant{Delay: 5 * time.Second, MaxRetries: 5}
	_, err := backoff.Wait(ctx, s, 0)
	if err == nil {
		t.Fatal("expected error on cancelled context")
	}
}

func TestWait_NoMoreRetries_ReturnsFalse(t *testing.T) {
	s := backoff.None{}
	ok, err := backoff.Wait(context.Background(), s, 0)
	if ok || err != nil {
		t.Fatalf("expected (false, nil), got (%v, %v)", ok, err)
	}
}
