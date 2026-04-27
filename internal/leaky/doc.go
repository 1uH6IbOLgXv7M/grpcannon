// Package leaky provides a leaky-bucket rate limiter.
//
// Unlike a token-bucket limiter that allows short bursts, a leaky bucket
// enforces a smooth, constant output rate. Incoming requests fill the
// bucket; the bucket drains at a fixed rate. When the bucket overflows
// the request is dropped with ErrDropped, giving the caller a clear
// signal to shed load rather than queue indefinitely.
//
// Typical usage:
//
//	bucket := leaky.New(100, 50) // capacity 100, drain 50 req/s
//	if err := bucket.Acquire(ctx); err != nil {
//		// shed or return 429
//	}
package leaky
