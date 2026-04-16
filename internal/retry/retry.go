// Package retry provides a configurable retry policy for gRPC invocations.
// It supports no-retry, fixed-count, and conditional retry strategies based
// on gRPC status codes.
package retry

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Policy describes how failed calls should be retried.
type Policy struct {
	// MaxAttempts is the total number of attempts (1 = no retry).
	MaxAttempts int
	// Delay is the fixed wait between attempts.
	Delay time.Duration
	// RetryOn is the set of gRPC codes that are retryable.
	// An empty set means retry on any error.
	RetryOn []codes.Code
}

// Default returns a Policy that never retries.
func Default() Policy {
	return Policy{MaxAttempts: 1}
}

// Do executes fn according to p, retrying on eligible errors.
// It returns the last error if all attempts are exhausted.
func Do(ctx context.Context, p Policy, fn func(ctx context.Context) error) error {
	max := p.MaxAttempts
	if max < 1 {
		max = 1
	}

	var last error
	for attempt := 0; attempt < max; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		last = fn(ctx)
		if last == nil {
			return nil
		}

		if !isRetryable(last, p.RetryOn) {
			return last
		}

		if attempt < max-1 && p.Delay > 0 {
			select {
			case <-time.After(p.Delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return last
}

// isRetryable reports whether err should trigger a retry given the allowed codes.
func isRetryable(err error, allowed []codes.Code) bool {
	if len(allowed) == 0 {
		return true
	}
	st, ok := status.FromError(err)
	if !ok {
		return false
	}
	for _, c := range allowed {
		if st.Code() == c {
			return true
		}
	}
	return false
}
