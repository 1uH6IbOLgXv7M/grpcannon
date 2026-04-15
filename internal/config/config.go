package config

import (
	"errors"
	"time"
)

// ConcurrencyProfile defines how concurrency ramps up during the test.
type ConcurrencyProfile string

const (
	ProfileFlat    ConcurrencyProfile = "flat"
	ProfileRampUp  ConcurrencyProfile = "ramp-up"
	ProfileStep    ConcurrencyProfile = "step"
)

// Config holds all runtime configuration for a grpcannon load test.
type Config struct {
	// Target is the gRPC server address (host:port).
	Target string `json:"target"`

	// Proto is the path to the .proto file describing the service.
	Proto string `json:"proto"`

	// Call is the fully-qualified method name, e.g. "pkg.Service/Method".
	Call string `json:"call"`

	// Data is the JSON-encoded request payload.
	Data string `json:"data"`

	// Concurrency is the number of concurrent workers.
	Concurrency int `json:"concurrency"`

	// TotalRequests is the total number of requests to send.
	TotalRequests int `json:"total_requests"`

	// Duration overrides TotalRequests when non-zero; test runs for this long.
	Duration time.Duration `json:"duration"`

	// Timeout is the per-request deadline.
	Timeout time.Duration `json:"timeout"`

	// Profile controls how concurrency is applied over time.
	Profile ConcurrencyProfile `json:"profile"`

	// Insecure disables TLS verification.
	Insecure bool `json:"insecure"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Concurrency:   10,
		TotalRequests: 200,
		Timeout:       5 * time.Second,
		Profile:       ProfileFlat,
		Insecure:      false,
	}
}

// Validate checks that the Config is complete and coherent.
func (c *Config) Validate() error {
	if c.Target == "" {
		return errors.New("target address is required")
	}
	if c.Call == "" {
		return errors.New("gRPC call (method) is required")
	}
	if c.Concurrency <= 0 {
		return errors.New("concurrency must be greater than zero")
	}
	if c.Duration == 0 && c.TotalRequests <= 0 {
		return errors.New("either duration or total_requests must be set")
	}
	if c.Timeout <= 0 {
		return errors.New("timeout must be greater than zero")
	}
	switch c.Profile {
	case ProfileFlat, ProfileRampUp, ProfileStep:
		// valid
	default:
		return errors.New("unknown concurrency profile: " + string(c.Profile))
	}
	return nil
}
