// Package profile provides concurrency profile primitives for grpcannon.
//
// A Profile is a sequence of Stages, each specifying a target worker count and
// a duration. The runner steps through stages in order, adjusting the worker
// pool size at each transition.
//
// Built-in helpers:
//
//	- Flat   – constant concurrency for the whole test
//	- Ramp   – linearly interpolate workers across N equal-length steps
//
// Example:
//
//	p := profile.Ramp(1, 50, 5, 10*time.Second)
//	// 5 stages of 10 s each: 1 → 13 → 25 → 38 → 50 workers
package profile
