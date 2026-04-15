package worker

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Result holds the outcome of a single gRPC call.
type Result struct {
	Duration time.Duration
	Err      error
}

// CallFunc is the function signature for executing a single gRPC call.
type CallFunc func(ctx context.Context) error

// Pool manages a fixed number of concurrent workers that execute gRPC calls.
type Pool struct {
	concurrency int
	total       int
	callFn      CallFunc
	Results     chan Result
	Completed   int64
}

// NewPool creates a new worker Pool.
func NewPool(concurrency, total int, fn CallFunc) *Pool {
	return &Pool{
		concurrency: concurrency,
		total:       total,
		callFn:      fn,
		Results:     make(chan Result, total),
	}
}

// Run dispatches work across the pool and closes Results when done.
func (p *Pool) Run(ctx context.Context) {
	work := make(chan struct{}, p.total)
	for i := 0; i < p.total; i++ {
		work <- struct{}{}
	}
	close(work)

	var wg sync.WaitGroup
	for i := 0; i < p.concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range work {
				if ctx.Err() != nil {
					return
				}
				start := time.Now()
				err := p.callFn(ctx)
				p.Results <- Result{
					Duration: time.Since(start),
					Err:      err,
				}
				atomic.AddInt64(&p.Completed, 1)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(p.Results)
	}()
}
