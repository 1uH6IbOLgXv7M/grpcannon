// Package retry implements a lightweight retry policy for use with gRPC
// invocations inside grpcannon.
//
// A [Policy] describes how many attempts to make, how long to wait between
// them, and which gRPC status codes are considered transient and therefore
// eligible for a retry.
//
// Usage:
//
//	p := retry.Policy{
//		MaxAttempts: 3,
//		Delay:       100 * time.Millisecond,
//		RetryOn:     []codes.Code{codes.Unavailable, codes.DeadlineExceeded},
//	}
//	err := retry.Do(ctx, p, func(ctx context.Context) error {
//		return invoker.Do(ctx, method, payload)
//	})
package retry
