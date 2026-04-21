package ewma

import (
	"math"
	"sync"
	"testing"
)

func TestNew_DefaultDecay_WhenAlphaInvalid(t *testing.T) {
	for _, alpha := range []float64{0, -1, 1.5} {
		e := New(alpha)
		if e.alpha != DefaultDecay {
			t.Errorf("alpha=%v: want DefaultDecay, got %v", alpha, e.alpha)
		}
	}
}

func TestNew_ValidAlpha_Preserved(t *testing.T) {
	e := New(0.3)
	if e.alpha != 0.3 {
		t.Fatalf("want 0.3, got %v", e.alpha)
	}
}

func TestValue_NoObservations_ReturnsZero(t *testing.T) {
	e := New(0.5)
	if e.Value() != 0 {
		t.Fatalf("want 0, got %v", e.Value())
	}
}

func TestAdd_FirstObservation_SetsValue(t *testing.T) {
	e := New(0.5)
	e.Add(42)
	if e.Value() != 42 {
		t.Fatalf("want 42, got %v", e.Value())
	}
}

func TestAdd_SecondObservation_Smoothed(t *testing.T) {
	e := New(0.5)
	e.Add(100)
	e.Add(0)
	// 0.5*0 + 0.5*100 = 50
	if e.Value() != 50 {
		t.Fatalf("want 50, got %v", e.Value())
	}
}

func TestAdd_MultipleObservations_ConvergesOnConstant(t *testing.T) {
	e := New(0.5)
	for i := 0; i < 100; i++ {
		e.Add(10)
	}
	if math.Abs(e.Value()-10) > 0.001 {
		t.Fatalf("want ~10, got %v", e.Value())
	}
}

func TestReset_ClearsValue(t *testing.T) {
	e := New(0.5)
	e.Add(99)
	e.Reset()
	if e.Value() != 0 {
		t.Fatalf("want 0 after reset, got %v", e.Value())
	}
	e.Add(7)
	if e.Value() != 7 {
		t.Fatalf("want 7 after reset+add, got %v", e.Value())
	}
}

func TestRate_RoundsToDecimals(t *testing.T) {
	e := New(0.5)
	e.Add(1)
	e.Add(2) // 0.5*2 + 0.5*1 = 1.5
	if e.Rate(0) != 2 {
		t.Fatalf("want 2, got %v", e.Rate(0))
	}
	if e.Rate(1) != 1.5 {
		t.Fatalf("want 1.5, got %v", e.Rate(1))
	}
}

func TestConcurrent_AddAndValue_NoRace(t *testing.T) {
	e := New(0.1)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(v float64) {
			defer wg.Done()
			e.Add(v)
			_ = e.Value()
		}(float64(i))
	}
	wg.Wait()
}
