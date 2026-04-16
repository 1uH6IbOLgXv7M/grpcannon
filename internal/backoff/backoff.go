// Package backoff provides simple retry backoff strategies for gRPC call failures.
package backoff

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// Strategy defines how long to wait before the next retry attempt.
type Strategy interface {
	// Next returns the duration to wait before attempt number n (0-indexed).
	// Returns false if no further retries should be attempted.
	Next(attempt int) (time.Duration, bool)
}

// None is a Strategy that never retries.
type None struct{}

func (None) Next(_ int) (time.Duration, bool) { return 0, false }

// Constant waits the same duration between every retry up to MaxRetries.
type Constant struct {
	Delay      time.Duration
	MaxRetries int
}

func (c Constant) Next(attempt int) (time.Duration, bool) {
	if attempt >= c.MaxRetries {
		return 0, false
	}
	return c.Delay, true
}

// Exponential implements truncated binary-exponential backoff with optional
// full jitter to spread retry storms.
type Exponential struct {
	BaseDelay  time.Duration
	MaxDelay   time.Duration
	MaxRetries int
	Jitter     bool
}

func (e Exponential) Next(attempt int) (time.Duration, bool) {
	if attempt >= e.MaxRetries {
		return 0, false
	}
	exp := math.Pow(2, float64(attempt))
	d := time.Duration(float64(e.BaseDelay) * exp)
	if d > e.MaxDelay || d <= 0 {
		d = e.MaxDelay
	}
	if e.Jitter {
		//nolint:gosec // non-cryptographic jitter is fine
		d = time.Duration(rand.Int63n(int64(d) + 1))
	}
	return d, true
}

// Wait blocks for the duration prescribed by s for the given attempt.
// Returns an error if ctx is cancelled before the wait completes.
// Returns (false, nil) when the strategy signals no more retries.
func Wait(ctx context.Context, s Strategy, attempt int) (bool, error) {
	d, ok := s.Next(attempt)
	if !ok {
		return false, nil
	}
	if d <= 0 {
		return true, ctx.Err()
	}
	select {
	case <-time.After(d):
		return true, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}
