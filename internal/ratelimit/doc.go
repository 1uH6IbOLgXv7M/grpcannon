// Package ratelimit implements a simple token-bucket rate limiter used by
// grpcannon to cap the number of outgoing gRPC requests per second.
//
// Usage:
//
//	// Allow up to 500 requests per second.
//	limiter := ratelimit.New(500)
//	defer limiter.Stop()
//
//	// Inside the hot loop:
//	if err := limiter.Wait(ctx); err != nil {
//		// context cancelled — stop sending.
//		return
//	}
//	// ... dispatch gRPC call ...
//
// Pass ratelimit.Unlimited (0) to disable throttling entirely.
package ratelimit
