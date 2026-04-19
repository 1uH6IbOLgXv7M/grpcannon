// Package cascade implements a consecutive-failure detector for use in
// load test pipelines. Unlike a circuit breaker it does not self-heal
// on a timer; callers must explicitly record a success to reset the state.
//
// Typical usage:
//
//	det := cascade.New(5)
//	if err := det.Allow(); err != nil {
//	    // skip request — too many consecutive failures
//	}
//	// … invoke RPC …
//	if rpcErr != nil {
//	    det.RecordFailure()
//	} else {
//	    det.RecordSuccess()
//	}
package cascade
