// Package worker provides a concurrency-controlled worker pool for executing
// gRPC calls during load tests.
//
// A Pool is created with a fixed concurrency level and a total request count.
// Callers supply a CallFunc that performs a single gRPC invocation; the pool
// dispatches exactly `total` calls across `concurrency` goroutines and streams
// each Result (duration + error) to the Results channel.
//
// Example usage:
//
//	pool := worker.NewPool(10, 1000, func(ctx context.Context) error {
//		// perform gRPC call
//		return nil
//	})
//	pool.Run(ctx)
//	for result := range pool.Results {
//		// collect latency / error data
//	}
package worker
