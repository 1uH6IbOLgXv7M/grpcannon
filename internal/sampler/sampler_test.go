package sampler

import (
	"fmt"
	"testing"
)

func TestNew_ZeroCapacity_NoSamplesRetained(t *testing.T) {
	s := New(0, 42)
	s.Add("svc/Method", []byte("data"))
	if got := len(s.Samples()); got != 0 {
		t.Fatalf("expected 0 samples, got %d", got)
	}
}

func TestNew_NegativeCapacity_TreatedAsZero(t *testing.T) {
	s := New(-5, 42)
	s.Add("svc/Method", []byte("data"))
	if got := len(s.Samples()); got != 0 {
		t.Fatalf("expected 0 samples, got %d", got)
	}
}

func TestAdd_BelowCapacity_AllRetained(t *testing.T) {
	s := New(10, 1)
	for i := 0; i < 7; i++ {
		s.Add("svc/M", []byte(fmt.Sprintf("payload-%d", i)))
	}
	if got := len(s.Samples()); got != 7 {
		t.Fatalf("expected 7 samples, got %d", got)
	}
}

func TestAdd_ExceedsCapacity_ReservoirSizeCapped(t *testing.T) {
	capacity := 5
	s := New(capacity, 99)
	for i := 0; i < 100; i++ {
		s.Add("svc/M", []byte(fmt.Sprintf("p%d", i)))
	}
	if got := len(s.Samples()); got != capacity {
		t.Fatalf("expected %d samples, got %d", capacity, got)
	}
}

func TestCount_TracksAllOffered(t *testing.T) {
	s := New(3, 7)
	for i := 0; i < 50; i++ {
		s.Add("svc/M", []byte("x"))
	}
	if got := s.Count(); got != 50 {
		t.Fatalf("expected count 50, got %d", got)
	}
}

func TestSamples_ReturnsCopy_MutationDoesNotAffectInternal(t *testing.T) {
	s := New(5, 11)
	s.Add("svc/M", []byte("hello"))
	snap := s.Samples()
	snap[0].Method = "tampered"

	original := s.Samples()
	if original[0].Method == "tampered" {
		t.Fatal("mutation of snapshot affected internal state")
	}
}

func TestAdd_PayloadIsCopied(t *testing.T) {
	s := New(5, 22)
	payload := []byte("original")
	s.Add("svc/M", payload)

	// Mutate the original slice after adding.
	payload[0] = 'X'

	got := s.Samples()[0].Payload
	if got[0] == 'X' {
		t.Fatal("sampler did not copy payload: mutation affected stored sample")
	}
}

func TestAdd_ConcurrentSafe(t *testing.T) {
	s := New(50, 33)
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				s.Add("svc/M", []byte(fmt.Sprintf("%d-%d", id, j)))
			}
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	if s.Count() != 1000 {
		t.Fatalf("expected count 1000, got %d", s.Count())
	}
	if got := len(s.Samples()); got > 50 {
		t.Fatalf("reservoir exceeded capacity: got %d", got)
	}
}
