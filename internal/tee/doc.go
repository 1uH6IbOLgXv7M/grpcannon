// Package tee provides a fan-out multiplexer for metrics.Snapshot values.
//
// A Tee holds a list of Sink implementations and delivers each incoming
// snapshot to all of them concurrently. Typical sinks include the terminal
// reporter, a JSON file writer, and the progress bar printer.
//
// Usage:
//
//	t := tee.New(reporter, fileWriter)
//	t.Run(ctx, snapshotCh)
package tee
