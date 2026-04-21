// Package scatter provides a generic key-based router that dispatches
// values to one of several named sinks. A routing key is extracted from
// each value by a caller-supplied function; if no sink is registered for
// the key the value is forwarded to an optional fallback sink.
//
// Typical use: fan-out snapshot streams by worker tag so that per-worker
// metric recorders can be updated independently without contention.
package scatter
