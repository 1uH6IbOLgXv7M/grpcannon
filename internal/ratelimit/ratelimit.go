// Package ratelimit provides a token-bucket rate limiter for controlling
// the request throughput during a load test run.
package ratelimit

import (
	"context"
	"time"
)

// Limiter controls the rate at which requests are dispatched.
type Limiter struct {
	ticker *time.Ticker
	done   chan struct{}
}

// Unlimited is a sentinel value meaning no rate limiting is applied.
const Unlimited = 0

// New creates a Limiter that allows rps requests per second.
// If rps is Unlimited (0) the limiter never blocks.
func New(rps int) *Limiter {
	if rps == Unlimited {
		return &Limiter{}
	}
	interval := time.Second / time.Duration(rps)
	return &Limiter{
		ticker: time.NewTicker(interval),
		done:   make(chan struct{}),
	}
}

// Wait blocks until the next token is available or ctx is cancelled.
// Returns ctx.Err() if the context is done before a token arrives.
func (l *Limiter) Wait(ctx context.Context) error {
	if l.ticker == nil {
		// Unlimited — return immediately.
		return ctx.Err()
	}
	select {
	case <-l.ticker.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop releases resources held by the Limiter.
func (l *Limiter) Stop() {
	if l.ticker != nil {
		l.ticker.Stop()
	}
}
