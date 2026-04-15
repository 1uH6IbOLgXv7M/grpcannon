package config

import (
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	c := Default()
	if c.Concurrency != 10 {
		t.Errorf("expected default concurrency 10, got %d", c.Concurrency)
	}
	if c.TotalRequests != 200 {
		t.Errorf("expected default total_requests 200, got %d", c.TotalRequests)
	}
	if c.Timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", c.Timeout)
	}
	if c.Profile != ProfileFlat {
		t.Errorf("expected default profile 'flat', got %s", c.Profile)
	}
	if c.Insecure {
		t.Error("expected insecure to be false by default")
	}
}

func TestValidate_Valid(t *testing.T) {
	c := Default()
	c.Target = "localhost:50051"
	c.Call = "example.Greeter/SayHello"
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_MissingTarget(t *testing.T) {
	c := Default()
	c.Call = "example.Greeter/SayHello"
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing target")
	}
}

func TestValidate_MissingCall(t *testing.T) {
	c := Default()
	c.Target = "localhost:50051"
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing call")
	}
}

func TestValidate_InvalidConcurrency(t *testing.T) {
	c := Default()
	c.Target = "localhost:50051"
	c.Call = "example.Greeter/SayHello"
	c.Concurrency = 0
	if err := c.Validate(); err == nil {
		t.Error("expected error for zero concurrency")
	}
}

func TestValidate_DurationOverridesTotalRequests(t *testing.T) {
	c := Default()
	c.Target = "localhost:50051"
	c.Call = "example.Greeter/SayHello"
	c.TotalRequests = 0
	c.Duration = 10 * time.Second
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error when duration is set, got: %v", err)
	}
}

func TestValidate_UnknownProfile(t *testing.T) {
	c := Default()
	c.Target = "localhost:50051"
	c.Call = "example.Greeter/SayHello"
	c.Profile = "unknown"
	if err := c.Validate(); err == nil {
		t.Error("expected error for unknown profile")
	}
}
