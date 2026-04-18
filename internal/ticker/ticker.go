// Package ticker provides a configurable wall-clock ticker that emits
// ticks on a channel and can be stopped cleanly.
package ticker

import (
	"context"
	"time"
)

// Ticker emits ticks at a fixed interval until stopped or the context is done.
type Ticker struct {
	interval time.Duration
	ch       chan time.Time
}

// New creates a Ticker that emits at the given interval.
// Intervals <= 0 default to 1 second.
func New(interval time.Duration) *Ticker {
	if interval <= 0 {
		interval = time.Second
	}
	return &Ticker{
		interval: interval,
		ch:       make(chan time.Time, 1),
	}
}

// C returns the channel on which ticks are delivered.
func (t *Ticker) C() <-chan time.Time {
	return t.ch
}

// Run starts emitting ticks until ctx is cancelled. It closes the channel
// when it exits so consumers can range over C().
func (t *Ticker) Run(ctx context.Context) {
	defer close(t.ch)
	tk := time.NewTicker(t.interval)
	defer tk.Stop()
	for {
		select {
		case now := <-tk.C:
			select {
			case t.ch <- now:
			default:
			}
		case <-ctx.Done():
			return
		}
	}
}

// Interval returns the configured tick interval.
func (t *Ticker) Interval() time.Duration {
	return t.interval
}
