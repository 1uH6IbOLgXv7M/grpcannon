// Package jitter provides helpers for introducing controlled randomness into
// time.Duration values.
//
// Three strategies are available:
//
//   - Full: uniform jitter in [0, d)
//   - Equal: uniform jitter in [d/2, d)
//   - Deviation: jitter within ±factor of d
//
// All functions are safe for concurrent use.
package jitter
