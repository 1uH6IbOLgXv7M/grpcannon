// Package budget implements an error-budget gate for grpcannon load runs.
//
// A Budget is configured with a failure-ratio threshold (e.g. 0.05 for 5 %).
// Callers record the outcome of every request via Record, then check Allow
// before dispatching the next request. Once the cumulative failure ratio
// reaches the threshold Allow returns ErrExceeded, letting the runner
// abort or throttle the run early.
//
// The zero value is not usable; construct a Budget with New.
package budget
