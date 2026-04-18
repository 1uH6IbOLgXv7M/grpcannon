// Package debounce provides a Debouncer that coalesces rapid successive
// calls into a single execution after a configurable quiet period.
//
// Typical use-cases include rate-limiting reactive updates triggered by
// high-frequency events such as metric flushes or config reloads.
//
// Usage:
//
//	d := debounce.New(200*time.Millisecond, func() {
//		fmt.Println("fired")
//	})
//	d.Call() // resets the timer on every invocation
//	d.Flush() // fires immediately if a call is pending
//	d.Stop()  // cancels without firing
package debounce
