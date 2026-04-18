// Package adaptive implements a feedback-driven concurrency controller for
// grpcannon load tests.
//
// A Controller tracks the success/failure ratio of recent requests and
// periodically adjusts the target worker count up or down within configured
// bounds. Callers should invoke Record after every request and Adjust on a
// regular interval (e.g. once per second) to obtain the new concurrency level.
package adaptive
