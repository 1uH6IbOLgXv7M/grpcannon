// Package watchdog provides an error-rate monitor that cancels an in-progress
// load test when observed failures exceed a configurable fraction of total
// requests.
//
// Usage:
//
//	ctx, cancel := context.WithCancelCause(parent)
//	wd := watchdog.New(watchdog.Config{Threshold: 0.05, MinRequests: 20}, src)
//	go wd.Run(ctx, cancel)
package watchdog
