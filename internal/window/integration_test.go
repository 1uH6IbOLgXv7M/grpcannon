package window_test

import (
	"sync"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/window"
)

func TestConcurrent_AddAndCounts(t *testing.T) {
	w := window.New(time.Second, 10)
	const goroutines = 20
	const perG = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < perG; j++ {
				w.Add(j%5 == 0)
			}
		}(i)
	}
	wg.Wait()
	total, errors := w.Counts()
	expected := int64(goroutines * perG)
	if total != expected {
		t.Fatalf("expected total %d, got %d", expected, total)
	}
	expectedErr := int64(goroutines * (perG / 5))
	if errors != expectedErr {
		t.Fatalf("expected errors %d, got %d", expectedErr, errors)
	}
}
