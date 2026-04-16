package snapshot

import (
	"context"
	"time"
)

// Ticker drives periodic snapshot collection at a fixed interval.
type Ticker struct {
	interval  time.Duration
	collector *Collector
	ch        chan Snapshot
}

// NewTicker creates a Ticker that captures a snapshot every interval.
func NewTicker(c *Collector, interval time.Duration) *Ticker {
	return &Ticker{
		interval:  interval,
		collector: c,
		ch:        make(chan Snapshot, 8),
	}
}

// C returns the channel on which snapshots are delivered.
func (t *Ticker) C() <-chan Snapshot { return t.ch }

// Run starts the ticker loop and blocks until ctx is cancelled.
func (t *Ticker) Run(ctx context.Context) {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()
	defer close(t.ch)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s := t.collector.Capture()
			select {
			case t.ch <- s:
			default:
			}
		}
	}
}
