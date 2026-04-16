package budget_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/your-org/grpcannon/internal/budget"
)

// TestConcurrent_RecordAndAllow exercises Budget under concurrent access.
func TestConcurrent_RecordAndAllow(t *testing.T) {
	b := budget.New(0.9) // very lenient so Allow rarely trips
	const goroutines = 50
	const ops = 200

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < ops; j++ {
				if id%5 == 0 {
					b.Record(errors.New("err"))
				} else {
					b.Record(nil)
				}
				_ = b.Allow()
			}
		}(i)
	}
	wg.Wait()

	ratio := b.Ratio()
	if ratio < 0 || ratio > 1 {
		t.Fatalf("ratio out of range: %f", ratio)
	}
	expected := float64(goroutines/5) / float64(goroutines)
	// allow 2 % tolerance
	if ratio < expected-0.02 || ratio > expected+0.02 {
		t.Fatalf("ratio %f far from expected ~%f", ratio, expected)
	}
}
