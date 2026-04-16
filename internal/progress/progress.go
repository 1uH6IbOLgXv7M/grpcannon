// Package progress provides a simple terminal progress reporter
// that prints live request counts and error rates during a load test.
package progress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/bojand/grpcannon/internal/snapshot"
)

// Reporter periodically writes a single-line progress update to an io.Writer.
type Reporter struct {
	mu       sync.Mutex
	out      io.Writer
	collector *snapshot.Collector
	interval time.Duration
	stop     chan struct{}
	wg       sync.WaitGroup
}

// New creates a Reporter that reads from collector every interval.
// If out is nil, os.Stderr is used.
func New(out io.Writer, collector *snapshot.Collector, interval time.Duration) *Reporter {
	if out == nil {
		out = os.Stderr
	}
	if interval <= 0 {
		interval = time.Second
	}
	return &Reporter{
		out:      out,
		collector: collector,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start begins background reporting. Call Stop to halt.
func (r *Reporter) Start() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.print()
			case <-r.stop:
				return
			}
		}
	}()
}

// Stop halts background reporting and waits for the goroutine to exit.
func (r *Reporter) Stop() {
	close(r.stop)
	r.wg.Wait()
}

func (r *Reporter) print() {
	snap := r.collector.Latest()
	errRate := 0.0
	if snap.Total > 0 {
		errRate = float64(snap.Errors) / float64(snap.Total) * 100
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	fmt.Fprintf(r.out, "\rprogress: total=%-6d errors=%-6d err_rate=%5.1f%%  p50=%v p99=%v",
		snap.Total, snap.Errors, errRate,
		snap.P50, snap.P99,
	)
}
