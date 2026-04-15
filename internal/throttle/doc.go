// Package throttle implements a token-bucket rate limiter used to cap
// the number of gRPC requests dispatched per second during a load test.
//
// When the configured RPS is zero or negative the throttle is unlimited
// and every call to Wait returns immediately without blocking.
//
// Typical usage:
//
//	th := throttle.New(500) // 500 req/s
//	defer th.Stop()
//	if err := th.Wait(ctx); err != nil {
//		// context cancelled
//	}
package throttle
