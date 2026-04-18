package window

import (
	"testing"
	"time"
)

func BenchmarkAdd(b *testing.B) {
	w := New(time.Second, 10)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w.Add(false)
		}
	})
}

func BenchmarkCounts(b *testing.B) {
	w := New(time.Second, 10)
	for i := 0; i < 1000; i++ {
		w.Add(i%10 == 0)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Counts()
	}
}
