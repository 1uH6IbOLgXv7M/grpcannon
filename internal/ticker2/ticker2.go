// Package ticker2 provides a wall-clock ticker that emits the current time
// at a fixed interval and can be stopped cleanly.
package ticker2

import "time"

// Ticker wraps time.Ticker with a clean Stop/channel API.
type Ticker struct {
	t  *time.Ticker
	C  <-chan time.Time
}

// New creates a Ticker that fires every d. If d <= 0 it defaults to 1 second.
func New(d time.Duration) *Ticker {
	if d <= 0 {
		d = time.Second
	}
	tt := time.NewTicker(d)
	return &Ticker{t: tt, C: tt.C}
}

// Stop halts the ticker. Subsequent reads on C will block forever.
func (t *Ticker) Stop() {
	t.t.Stop()
}

// Reset changes the ticker interval. It stops the old ticker and starts a new
// one; callers must re-read from C after Reset returns.
func (t *Ticker) Reset(d time.Duration) {
	if d <= 0 {
		d = time.Second
	}
	t.t.Reset(d)
}
