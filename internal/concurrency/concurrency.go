// Package concurrency provides utilities for managing dynamic worker
// concurrency during a load test, applying a staged profile over time.
package concurrency

import (
	"context"
	"time"

	"github.com/nicklaw5/grpcannon/internal/profile"
)

// Controller adjusts the number of active workers according to a concurrency
// profile, signalling changes over a channel.
type Controller struct {
	stages []profile.Stage
	changes chan int
}

// New creates a Controller for the given profile stages.
func New(stages []profile.Stage) *Controller {
	return &Controller{
		stages: stages,
		changes: make(chan int, 1),
	}
}

// Changes returns a read-only channel that emits the desired worker count
// whenever the concurrency level should change.
func (c *Controller) Changes() <-chan int {
	return c.changes
}

// Run drives the controller through each stage, blocking until all stages
// complete or ctx is cancelled. The final worker count is emitted before Run
// returns so callers can always read at least one value.
func (c *Controller) Run(ctx context.Context) {
	defer close(c.changes)

	for _, s := range c.stages {
		select {
		case <-ctx.Done():
			return
		case c.changes <- s.Workers:
		}

		timer := time.NewTimer(s.Duration)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}
