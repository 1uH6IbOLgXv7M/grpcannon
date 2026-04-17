// Package fence provides a one-shot barrier (Fence) that goroutines can wait
// on until it is explicitly opened. Once opened the fence remains open and all
// subsequent Wait calls return immediately.
//
// A typical use-case is delaying worker goroutines until the load-test
// start signal has been confirmed (e.g. after a warm-up phase completes).
package fence
