// Package signal wraps OS signal handling for grpcannon's graceful shutdown
// path. It provides two helpers:
//
//   - NotifyContext: returns a context that is cancelled on SIGINT/SIGTERM,
//     following the standard library's signal.NotifyContext pattern.
//
//   - Trap: a blocking helper suitable for running in a dedicated goroutine;
//     it calls a cancel function when an interrupt signal arrives.
//
// Both helpers ensure that signal handlers are cleaned up after use so that
// default signal behaviour is restored once the handler fires.
package signal
