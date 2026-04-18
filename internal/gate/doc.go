// Package gate provides a concurrency gate — a lightweight semaphore that
// limits the number of goroutines allowed to proceed simultaneously.
//
// Unlike a raw channel semaphore, Gate integrates context cancellation and
// supports a Close operation that permanently stops new callers from entering,
// making it suitable for graceful-shutdown sequences in load-test runners.
package gate
