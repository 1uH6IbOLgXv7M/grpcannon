package pacemaker

import (
	"testing"
	"time"
)

func cfg(target, min, max float64) Config {
	return Config{
		TargetP99:  time.Duration(target) * time.Millisecond,
		MinRPS:     min,
		MaxRPS:     max,
		StepFactor: 0.10,
	}
}

func TestNew_DefaultsToMaxRPS(t *testing.T) {
	p := New(cfg(50, 1, 200))
	if p.Current() != 200 {
		t.Fatalf("expected 200, got %f", p.Current())
	}
}

func TestNew_NoMaxRPS_Defaults100(t *testing.T) {
	p := New(cfg(50, 0, 0))
	if p.Current() != 100 {
		t.Fatalf("expected 100, got %f", p.Current())
	}
}

func TestAdjust_NoTarget_ReturnsCurrent(t *testing.T) {
	p := New(Config{MaxRPS: 50})
	got := p.Adjust(200 * time.Millisecond)
	if got != p.Current() {
		t.Fatalf("expected unchanged current")
	}
}

func TestAdjust_HighLatency_ReducesRPS(t *testing.T) {
	p := New(cfg(50, 1, 200))
	before := p.Current()
	after := p.Adjust(200 * time.Millisecond) // 4× target
	if after >= before {
		t.Fatalf("expected rate to decrease, got %f -> %f", before, after)
	}
}

func TestAdjust_LowLatency_IncreasesRPS(t *testing.T) {
	p := New(cfg(50, 1, 500))
	// seed at a lower value so there is room to grow
	p.current = 100
	after := p.Adjust(10 * time.Millisecond) // well below target
	if after <= 100 {
		t.Fatalf("expected rate to increase, got %f", after)
	}
}

func TestAdjust_ClampsToMin(t *testing.T) {
	p := New(cfg(50, 10, 200))
	p.current = 11
	// extremely high latency should not push below MinRPS
	for i := 0; i < 100; i++ {
		p.Adjust(10 * time.Second)
	}
	if p.Current() < 10 {
		t.Fatalf("rate dropped below MinRPS: %f", p.Current())
	}
}

func TestAdjust_ClampsToMax(t *testing.T) {
	p := New(cfg(50, 1, 200))
	p.current = 190
	for i := 0; i < 100; i++ {
		p.Adjust(1 * time.Millisecond)
	}
	if p.Current() > 200 {
		t.Fatalf("rate exceeded MaxRPS: %f", p.Current())
	}
}

func TestAdjust_NearTarget_NoChange(t *testing.T) {
	p := New(cfg(50, 1, 200))
	p.current = 150
	// p99 == 95 % of target: inside the dead-band [0.9, 1.0]
	got := p.Adjust(47 * time.Millisecond)
	if got != 150 {
		t.Fatalf("expected no change, got %f", got)
	}
}
