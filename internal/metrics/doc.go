// Package metrics implements latency recording and statistical aggregation
// for grpcannon load test runs.
//
// # Overview
//
// A [Recorder] collects individual RPC latency samples and error counts
// produced by the worker pool during a run. Once the run completes, call
// [Recorder.Snapshot] to obtain an immutable [Summary] containing:
//
//   - Total request count and error count
//   - Min / Mean / Max latency
//   - P50 / P95 / P99 latency percentiles
//
// The Recorder is safe for concurrent use by multiple goroutines.
//
// # Example
//
//	rec := metrics.NewRecorder()
//	// … pass rec.Record to each worker …
//	fmt.Println(rec.Snapshot())
package metrics
