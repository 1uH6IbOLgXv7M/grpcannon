// Package window implements a sliding-window counter that partitions a
// rolling time interval into fixed-size slots. Each slot accumulates
// request and error counts; slots older than the full window duration
// are discarded on the next read, giving an approximate rolling view
// of recent traffic suitable for error-rate guards and adaptive
// concurrency decisions inside grpcannon.
package window
