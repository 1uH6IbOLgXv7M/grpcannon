// Package cooldown provides a concurrency-safe cooldown timer that gates
// repeated actions to no more than once per configured interval.
//
// Typical use: suppress log spam, limit circuit-breaker probes, or throttle
// adaptive-concurrency adjustments to a human-readable cadence.
package cooldown
