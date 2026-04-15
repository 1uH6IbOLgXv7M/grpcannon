package profile_test

import (
	"testing"
	"time"

	"github.com/user/grpcannon/internal/profile"
)

func TestFlat_SingleStage(t *testing.T) {
	p := profile.Flat(10, 30*time.Second)
	if len(p.Stages) != 1 {
		t.Fatalf("expected 1 stage, got %d", len(p.Stages))
	}
	if p.Stages[0].Workers != 10 {
		t.Errorf("expected 10 workers, got %d", p.Stages[0].Workers)
	}
	if p.Stages[0].Duration != 30*time.Second {
		t.Errorf("unexpected duration: %v", p.Stages[0].Duration)
	}
}

func TestFlat_TotalDuration(t *testing.T) {
	p := profile.Flat(5, 20*time.Second)
	if p.TotalDuration() != 20*time.Second {
		t.Errorf("expected 20s, got %v", p.TotalDuration())
	}
}

func TestRamp_StageCount(t *testing.T) {
	p := profile.Ramp(1, 10, 5, 5*time.Second)
	if len(p.Stages) != 5 {
		t.Fatalf("expected 5 stages, got %d", len(p.Stages))
	}
}

func TestRamp_FirstAndLastWorkers(t *testing.T) {
	p := profile.Ramp(2, 20, 5, 10*time.Second)
	first := p.Stages[0].Workers
	last := p.Stages[len(p.Stages)-1].Workers
	if first != 2 {
		t.Errorf("expected first stage workers=2, got %d", first)
	}
	if last != 20 {
		t.Errorf("expected last stage workers=20, got %d", last)
	}
}

func TestRamp_MinSteps(t *testing.T) {
	// steps < 2 should be clamped to 2
	p := profile.Ramp(1, 10, 1, 5*time.Second)
	if len(p.Stages) != 2 {
		t.Errorf("expected 2 stages after clamp, got %d", len(p.Stages))
	}
}

func TestValidate_Valid(t *testing.T) {
	p := profile.Flat(4, 10*time.Second)
	if err := p.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_NoStages(t *testing.T) {
	p := &profile.Profile{}
	if err := p.Validate(); err == nil {
		t.Error("expected error for empty stages")
	}
}

func TestValidate_ZeroWorkers(t *testing.T) {
	p := &profile.Profile{
		Stages: []profile.Stage{{Workers: 0, Duration: 5 * time.Second}},
	}
	if err := p.Validate(); err == nil {
		t.Error("expected error for zero workers")
	}
}

func TestValidate_ZeroDuration(t *testing.T) {
	p := &profile.Profile{
		Stages: []profile.Stage{{Workers: 5, Duration: 0}},
	}
	if err := p.Validate(); err == nil {
		t.Error("expected error for zero duration")
	}
}

func TestTotalDuration_MultiStage(t *testing.T) {
	p := profile.Ramp(1, 8, 4, 10*time.Second)
	expected := 4 * 10 * time.Second
	if p.TotalDuration() != expected {
		t.Errorf("expected %v, got %v", expected, p.TotalDuration())
	}
}
