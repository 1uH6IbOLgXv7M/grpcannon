// Package concurrency implements a stage-driven concurrency controller for
// grpcannon load tests.
//
// A Controller accepts a slice of profile.Stage values and emits the desired
// worker count over a channel as each stage begins. Callers read from Changes()
// and resize their worker pool accordingly.
//
// Example:
//
//	stages := profile.Flat(50, 30*time.Second)
//	ctrl := concurrency.New(stages)
//	go ctrl.Run(ctx)
//	for workers := range ctrl.Changes() {
//		pool.Resize(workers)
//	}
package concurrency
