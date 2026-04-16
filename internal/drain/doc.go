// Package drain provides a Drainer type that enables graceful shutdown
// of a load test run. Callers call Acquire before dispatching each
// request and Release when the request completes. Once the test
// duration elapses, the runner calls Drain to block until all
// in-flight requests finish or a deadline is exceeded, ensuring
// metrics are fully recorded before the process exits.
package drain
