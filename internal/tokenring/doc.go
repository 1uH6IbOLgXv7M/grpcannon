// Package tokenring provides a thread-safe round-robin token ring for
// distributing integer slot identifiers evenly across concurrent workers.
//
// Typical usage in grpcannon is to assign each backend endpoint an integer
// index and call Next() before every outbound RPC so that load is spread
// uniformly without any single worker monopolising a slot.
package tokenring
