// Package signal provides graceful shutdown handling for grpcannon.
// It listens for OS interrupt signals and cancels a context so that
// in-flight requests can drain before the process exits.
package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// NotifyContext returns a derived context that is cancelled when the process
// receives SIGINT or SIGTERM. The returned stop function deregisters the
// signal handler and should be deferred by the caller.
func NotifyContext(parent context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
}

// Trap blocks until one of SIGINT or SIGTERM is received, then calls cancel.
// It is intended to be run in a dedicated goroutine alongside the main load
// test loop.
func Trap(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(ch)
	<-ch
	cancel()
}
