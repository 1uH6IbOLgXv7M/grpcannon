package jitter

import (
	"testing"
	"time"
)

func BenchmarkFull(b *testing.B) {
	d := 500 * time.Millisecond
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Full(d)
	}
}

func BenchmarkEqual(b *testing.B) {
	d := 500 * time.Millisecond
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Equal(d)
	}
}

func BenchmarkDeviation(b *testing.B) {
	d := 500 * time.Millisecond
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Deviation(d, 0.1)
	}
}
