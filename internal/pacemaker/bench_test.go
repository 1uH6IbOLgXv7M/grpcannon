package pacemaker

import (
	"testing"
	"time"
)

func BenchmarkAdjust(b *testing.B) {
	p := New(Config{
		TargetP99:  50 * time.Millisecond,
		MinRPS:     1,
		MaxRPS:     10000,
		StepFactor: 0.05,
	})
	latencies := []time.Duration{
		10 * time.Millisecond,
		45 * time.Millisecond,
		55 * time.Millisecond,
		200 * time.Millisecond,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Adjust(latencies[i%len(latencies)])
	}
}

func BenchmarkCurrent(b *testing.B) {
	p := New(Config{MaxRPS: 500})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Current()
	}
}
