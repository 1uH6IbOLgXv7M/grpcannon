// Package counter provides a lightweight, thread-safe atomic counter for
// tracking total requests and errors during a load test run.
//
// Counter is safe for concurrent use by multiple goroutines and is designed
// for minimal contention in high-throughput scenarios.
//
// Example:
//
//	c := counter.New()
//	c.IncTotal()
//	c.IncErrors()
//	fmt.Println(c.ErrorRate()) // 1.0
package counter
