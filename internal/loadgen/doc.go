// Package loadgen drives a configurable number of concurrent gRPC workers
// according to a stream of concurrency-level stages.
//
// A stage channel is typically produced by internal/profile (Flat or Ramp)
// and consumed by Run, which grows or shrinks the worker pool whenever a new
// value arrives.  Each worker calls the supplied RequestFunc in a tight loop,
// honouring an optional RPS throttle and recording every observation into a
// metrics.Recorder.
//
// Example usage:
//
//	stages := profile.Flat(profile.Stage{Workers: 10, Duration: 30 * time.Second})
//	rec   := metrics.NewRecorder()
//	err   := loadgen.Run(ctx, loadgen.Config{
//		Stages:   stages,
//		RPS:      500,
//		Recorder: rec,
//		Fn:       myGRPCCall,
//	})
package loadgen
