// Package ticker provides a wall-clock ticker used to drive periodic
// operations such as progress reporting and snapshot collection inside
// grpcannon. It wraps time.Ticker with context-aware lifecycle management
// and a non-blocking send so slow consumers never stall the tick loop.
package ticker
