// Package estop implements an emergency-stop (e-stop) guard for load
// generation runs.
//
// An EStop monitors the ratio of failures to total observations. Once that
// ratio exceeds the configured threshold the stop is tripped and every
// subsequent call to Allow returns ErrTripped, halting further dispatch.
//
// The stop can be Reset between runs to clear both the tripped flag and the
// internal counters, making it safe to reuse across multiple test phases.
//
// Typical usage:
//
//	es := estop.New(0.10) // trip if >10 % errors
//	for _, req := range requests {
//		if err := es.Allow(); err != nil {
//			break
//		}
//		if err := dispatch(req); err != nil {
//			es.RecordFailure()
//		} else {
//			es.RecordSuccess()
//		}
//	}
package estop
