// Package shed provides a simple load-shedding primitive for grpcannon.
//
// It tracks the number of concurrent in-flight gRPC requests and rejects
// new acquisitions once a configured ceiling is reached, preventing runaway
// concurrency from overwhelming the target service during load tests.
package shed
