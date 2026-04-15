// Package profile defines concurrency profiles for grpcannon load tests.
// A profile describes how worker concurrency ramps up, holds, and ramps down
// over the duration of a test run.
package profile

import (
	"errors"
	"time"
)

// Stage represents a single step in a concurrency profile.
type Stage struct {
	// Workers is the target number of concurrent workers for this stage.
	Workers int
	// Duration is how long this stage lasts.
	Duration time.Duration
}

// Profile is an ordered sequence of stages that describe how concurrency
// changes over the lifetime of a load test.
type Profile struct {
	Stages []Stage
}

// Validate returns an error if the profile is not usable.
func (p *Profile) Validate() error {
	if len(p.Stages) == 0 {
		return errors.New("profile must have at least one stage")
	}
	for i, s := range p.Stages {
		if s.Workers <= 0 {
			return fmt.Errorf("stage %d: workers must be > 0", i)
		}
		if s.Duration <= 0 {
			return fmt.Errorf("stage %d: duration must be > 0", i)
		}
	}
	return nil
}

// TotalDuration returns the sum of all stage durations.
func (p *Profile) TotalDuration() time.Duration {
	var total time.Duration
	for _, s := range p.Stages {
		total += s.Duration
	}
	return total
}

// Flat returns a Profile with a single constant-concurrency stage.
func Flat(workers int, duration time.Duration) *Profile {
	return &Profile{
		Stages: []Stage{
			{Workers: workers, Duration: duration},
		},
	}
}

// Ramp returns a Profile that linearly steps from startWorkers to endWorkers
// across steps stages, each lasting stepDuration.
func Ramp(startWorkers, endWorkers, steps int, stepDuration time.Duration) *Profile {
	if steps < 2 {
		steps = 2
	}
	stages := make([]Stage, steps)
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		w := startWorkers + int(t*float64(endWorkers-startWorkers))
		if w < 1 {
			w = 1
		}
		stages[i] = Stage{Workers: w, Duration: stepDuration}
	}
	return &Profile{Stages: stages}
}
