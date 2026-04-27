package tokenring

import (
	"sync"
	"testing"
)

func TestNew_SizeOne_AlwaysReturnsZero(t *testing.T) {
	r := New(1)
	for i := 0; i < 5; i++ {
		v, err := r.Next()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != 0 {
			t.Fatalf("want 0, got %d", v)
		}
	}
}

func TestNew_ClampsBadSize(t *testing.T) {
	r := New(0)
	if r.Len() != 1 {
		t.Fatalf("want len 1 after clamp, got %d", r.Len())
	}
}

func TestNext_RoundRobin(t *testing.T) {
	r := New(3)
	want := []int{0, 1, 2, 0, 1, 2}
	for i, w := range want {
		v, err := r.Next()
		if err != nil {
			t.Fatalf("step %d: unexpected error: %v", i, err)
		}
		if v != w {
			t.Fatalf("step %d: want %d, got %d", i, w, v)
		}
	}
}

func TestNext_EmptyRing_ReturnsErrEmpty(t *testing.T) {
	r := New(1)
	r.Reset(0)
	_, err := r.Next()
	if err != ErrEmpty {
		t.Fatalf("want ErrEmpty, got %v", err)
	}
}

func TestReset_ChangesSize(t *testing.T) {
	r := New(4)
	r.Reset(2)
	if r.Len() != 2 {
		t.Fatalf("want len 2 after reset, got %d", r.Len())
	}
	v, err := r.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 0 {
		t.Fatalf("want 0 after reset, got %d", v)
	}
}

func TestConcurrent_NeverExceedsSlotRange(t *testing.T) {
	const size = 5
	const goroutines = 20
	const calls = 200

	r := New(size)
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < calls; i++ {
				v, err := r.Next()
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if v < 0 || v >= size {
					t.Errorf("token %d out of range [0, %d)", v, size)
				}
			}
		}()
	}
	wg.Wait()
}
