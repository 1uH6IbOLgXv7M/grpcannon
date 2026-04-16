// Package circuit implements a simple circuit breaker that opens after a
// configurable number of consecutive failures and resets after a cooldown.
package circuit

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit breaker is open.
var ErrOpen = errors.New("circuit breaker open")

// State represents the current state of the breaker.
type State int

const (
	StateClosed State = iota
	StateOpen
)

// Breaker is a simple circuit breaker.
type Breaker struct {
	mu          sync.Mutex
	maxFailures int
	cooldown    time.Duration
	failures    int
	openedAt    time.Time
	state       State
	now         func() time.Time
}

// New returns a Breaker that opens after maxFailures consecutive failures
// and attempts to close again after cooldown.
func New(maxFailures int, cooldown time.Duration) *Breaker {
	if maxFailures < 1 {
		maxFailures = 1
	}
	return &Breaker{
		maxFailures: maxFailures,
		cooldown:    cooldown,
		now:         time.Now,
	}
}

// Allow returns nil if the call should proceed, or ErrOpen if the circuit is open.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == StateOpen {
		if b.now().Sub(b.openedAt) >= b.cooldown {
			b.state = StateClosed
			b.failures = 0
		} else {
			return ErrOpen
		}
	}
	return nil
}

// RecordSuccess resets the consecutive failure counter.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
}

// RecordFailure increments the failure counter and opens the circuit if the
// threshold is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.maxFailures {
		b.state = StateOpen
		b.openedAt = b.now()
	}
}

// State returns the current state of the breaker.
func (b *Breaker) CurrentState() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
