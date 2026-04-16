// Package deadline provides per-request timeout enforcement.
package deadline

import (
	"context"
	"errors"
	"time"
)

// ErrExceeded is returned when a request exceeds its allowed deadline.
var ErrExceeded = errors.New("deadline exceeded")

// Enforcer wraps a context with a per-request timeout.
type Enforcer struct {
	timeout time.Duration
}

// New returns an Enforcer that cancels derived contexts after timeout.
// A zero or negative timeout means no per-request deadline is applied.
func New(timeout time.Duration) *Enforcer {
	return &Enforcer{timeout: timeout}
}

// Wrap returns a child context and cancel func. If the Enforcer has a
// positive timeout the context will expire after that duration;
// otherwise the parent context is returned unchanged.
func (e *Enforcer) Wrap(parent context.Context) (context.Context, context.CancelFunc) {
	if e.timeout <= 0 {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, e.timeout)
}

// IsExceeded reports whether err represents a deadline/timeout condition
// originating from this package or from context itself.
func IsExceeded(err error) bool {
	if errors.Is(err, ErrExceeded) {
		return true
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	return false
}
