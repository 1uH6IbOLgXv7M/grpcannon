// Package probe implements a lightweight health-check prober that periodically
// sends a single gRPC request and reports whether the target is reachable.
// It is used by the load generator to gate traffic before the main run begins
// and to detect when a target has become unhealthy mid-test.
package probe

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// Status is the result of a single probe attempt.
type Status struct {
	// Healthy is true when the probe call returned without error.
	Healthy bool
	// Latency is the round-trip time of the probe call.
	Latency time.Duration
	// Err holds the error returned by the probe call, or nil.
	Err error
	// At is the wall-clock time at which the probe completed.
	At time.Time
}

// Fn is the function called by the prober on each tick.
// Implementations should make a single, lightweight gRPC call and return
// any error. The call must respect the supplied context.
type Fn func(ctx context.Context, conn *grpc.ClientConn) error

// Config holds the configuration for a Prober.
type Config struct {
	// Interval between successive probe attempts. Defaults to 1 s.
	Interval time.Duration
	// Timeout for a single probe call. Defaults to 5 s.
	Timeout time.Duration
	// Threshold is the number of consecutive failures before the prober
	// reports the target as unhealthy. Defaults to 1.
	Threshold int
}

func (c *Config) defaults() {
	if c.Interval <= 0 {
		c.Interval = time.Second
	}
	if c.Timeout <= 0 {
		c.Timeout = 5 * time.Second
	}
	if c.Threshold < 1 {
		c.Threshold = 1
	}
}

// Prober periodically calls a Fn and exposes the latest Status.
type Prober struct {
	cfg  Config
	conn *grpc.ClientConn
	fn   Fn

	mu       sync.RWMutex
	latest   Status
	failures int
}

// New creates a Prober. conn must not be nil; fn must not be nil.
func New(conn *grpc.ClientConn, fn Fn, cfg Config) *Prober {
	if conn == nil {
		panic("probe: conn must not be nil")
	}
	if fn == nil {
		panic("probe: fn must not be nil")
	}
	cfg.defaults()
	return &Prober{cfg: cfg, conn: conn, fn: fn}
}

// Latest returns the most recent probe Status. Before the first probe
// completes it returns a zero-value Status with Healthy == false.
func (p *Prober) Latest() Status {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.latest
}

// Healthy reports whether the most recent probe succeeded.
func (p *Prober) Healthy() bool {
	return p.Latest().Healthy
}

// Run starts the probe loop. It blocks until ctx is cancelled, then returns
// nil. Callers should run it in a goroutine.
func (p *Prober) Run(ctx context.Context) error {
	ticker := time.NewTicker(p.cfg.Interval)
	defer ticker.Stop()

	// Probe immediately before waiting for the first tick.
	p.probe(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			p.probe(ctx)
		}
	}
}

// probe executes a single probe attempt and updates the internal state.
func (p *Prober) probe(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, p.cfg.Timeout)
	defer cancel()

	start := time.Now()
	err := p.fn(ctx, p.conn)
	elapsed := time.Since(start)

	p.mu.Lock()
	defer p.mu.Unlock()

	if err != nil {
		p.failures++
		healthy := p.failures < p.cfg.Threshold
		p.latest = Status{
			Healthy: healthy,
			Latency: elapsed,
			Err:     fmt.Errorf("probe: %w", err),
			At:      time.Now(),
		}
		return
	}

	p.failures = 0
	p.latest = Status{
		Healthy: true,
		Latency: elapsed,
		At:      time.Now(),
	}
}
