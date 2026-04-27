package slope

import (
	"math"
	"testing"
)

func TestNew_ClampsBadWindow(t *testing.T) {
	e := New(0)
	if e.window != 2 {
		t.Fatalf("expected window=2, got %d", e.window)
	}
}

func TestNew_NegativeWindow_ClampsToTwo(t *testing.T) {
	e := New(-5)
	if e.window != 2 {
		t.Fatalf("expected window=2, got %d", e.window)
	}
}

func TestSlope_NoSamples_ReturnsZero(t *testing.T) {
	e := New(4)
	if s := e.Slope(); s != 0 {
		t.Fatalf("expected 0, got %f", s)
	}
}

func TestSlope_OneSample_ReturnsZero(t *testing.T) {
	e := New(4)
	e.Add(0, 10)
	if s := e.Slope(); s != 0 {
		t.Fatalf("expected 0, got %f", s)
	}
}

func TestSlope_PerfectlyLinearIncreasing(t *testing.T) {
	e := New(5)
	// y = 2x + 1  =>  slope should be exactly 2
	for i := 0; i < 5; i++ {
		e.Add(float64(i), float64(2*i+1))
	}
	got := e.Slope()
	if math.Abs(got-2.0) > 1e-9 {
		t.Fatalf("expected slope=2.0, got %f", got)
	}
}

func TestSlope_PerfectlyLinearDecreasing(t *testing.T) {
	e := New(4)
	// y = -3x
	for i := 0; i < 4; i++ {
		e.Add(float64(i), float64(-3*i))
	}
	got := e.Slope()
	if math.Abs(got-(-3.0)) > 1e-9 {
		t.Fatalf("expected slope=-3.0, got %f", got)
	}
}

func TestSlope_ConstantSeries_ReturnsZero(t *testing.T) {
	e := New(6)
	for i := 0; i < 6; i++ {
		e.Add(float64(i), 42)
	}
	got := e.Slope()
	if math.Abs(got) > 1e-9 {
		t.Fatalf("expected slope≈0, got %f", got)
	}
}

func TestSlope_WindowWrapsCorrectly(t *testing.T) {
	// Window of 3; add 6 samples so the buffer wraps twice.
	// Last 3 samples: (3,6),(4,8),(5,10)  => slope = 2
	e := New(3)
	for i := 0; i < 6; i++ {
		e.Add(float64(i), float64(2*i))
	}
	got := e.Slope()
	if math.Abs(got-2.0) > 1e-9 {
		t.Fatalf("expected slope=2.0 after wrap, got %f", got)
	}
}

func TestCount_BelowWindow(t *testing.T) {
	e := New(10)
	e.Add(0, 1)
	e.Add(1, 2)
	if c := e.Count(); c != 2 {
		t.Fatalf("expected count=2, got %d", c)
	}
}

func TestCount_AtWindow_ReturnsWindow(t *testing.T) {
	e := New(4)
	for i := 0; i < 6; i++ {
		e.Add(float64(i), float64(i))
	}
	if c := e.Count(); c != 4 {
		t.Fatalf("expected count=4, got %d", c)
	}
}
