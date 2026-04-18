// Package pacemaker implements a feedback-driven rate controller that adjusts
// the target requests-per-second based on observed p99 latency.
//
// Usage:
//
//	pm := pacemaker.New(pacemaker.Config{
//		TargetP99:  50 * time.Millisecond,
//		MinRPS:     10,
//		MaxRPS:     500,
//	})
//
//	// after each measurement window:
//	newRPS := pm.Adjust(measuredP99)
package pacemaker
