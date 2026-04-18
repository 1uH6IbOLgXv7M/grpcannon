// Package jitter provides utilities for adding randomised jitter to durations.
// Jitter is useful when spreading retry or backoff delays across many concurrent
// workers to avoid thundering-herd scenarios.
package jitter

import (
	"math/rand"
	"time"
)

// Source is a function that returns a pseudo-random float64 in [0,1).
type Source func() float64

// defaultSource uses the package-level rand which is safe for concurrent use
// since Go 1.20.
var defaultSource Source = rand.Float64

// Full returns a duration in [0, d).
func Full(d time.Duration) time.Duration {
	return full(d, defaultSource)
}

func full(d time.Duration, src Source) time.Duration {
	if d <= 0 {
		return 0
	}
	return time.Duration(src() * float64(d))
}

// Equal returns a duration in [d/2, d).
func Equal(d time.Duration) time.Duration {
	return equal(d, defaultSource)
}

func equal(d time.Duration, src Source) time.Duration {
	if d <= 0 {
		return 0
	}
	half := d / 2
	return half + time.Duration(src()*float64(half))
}

// Deviation returns a duration within ±factor of d.
// factor must be in (0, 1]; values outside this range are clamped.
func Deviation(d time.Duration, factor float64) time.Duration {
	return deviation(d, factor, defaultSource)
}

func deviation(d time.Duration, factor float64, src Source) time.Duration {
	if d <= 0 {
		return 0
	}
	if factor <= 0 {
		return d
	}
	if factor > 1 {
		factor = 1
	}
	delta := float64(d) * factor
	min := float64(d) - delta
	return time.Duration(min + src()*2*delta)
}
