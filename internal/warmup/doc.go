// Package warmup implements a pre-test warm-up phase for grpcannon.
//
// Before the main load test begins, a configurable number of requests are
// sent with limited concurrency. This allows the target service, any
// intermediate proxies, and the Go runtime itself to reach a steady state,
// reducing the impact of cold-start latency on benchmark results.
//
// Usage:
//
//	cfg := warmup.Default()
//	errCount, err := warmup.Run(ctx, cfg, "/pkg.Service/Method", payload, inv)
package warmup
