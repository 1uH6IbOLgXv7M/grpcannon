package topology

import (
	"sync"
	"testing"
)

func TestNew_EmptyEndpoints_ReturnsError(t *testing.T) {
	_, err := New(nil)
	if err != ErrNoEndpoints {
		t.Fatalf("expected ErrNoEndpoints, got %v", err)
	}

	_, err = New([]string{})
	if err != ErrNoEndpoints {
		t.Fatalf("expected ErrNoEndpoints for empty slice, got %v", err)
	}
}

func TestNew_ValidEndpoints_NoError(t *testing.T) {
	topo, err := New([]string{"localhost:50051"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if topo.Len() != 1 {
		t.Fatalf("expected Len 1, got %d", topo.Len())
	}
}

func TestNext_SingleEndpoint_AlwaysReturnsSame(t *testing.T) {
	topo, _ := New([]string{"host:1234"})
	for i := 0; i < 10; i++ {
		if got := topo.Next(); got != "host:1234" {
			t.Fatalf("iteration %d: expected host:1234, got %s", i, got)
		}
	}
}

func TestNext_MultipleEndpoints_RoundRobin(t *testing.T) {
	endpoints := []string{"a:1", "b:2", "c:3"}
	topo, _ := New(endpoints)

	for round := 0; round < 2; round++ {
		for i, want := range endpoints {
			got := topo.Next()
			if got != want {
				t.Fatalf("round %d, step %d: expected %s, got %s", round, i, want, got)
			}
		}
	}
}

func TestEndpoints_ReturnsCopy(t *testing.T) {
	original := []string{"x:9", "y:10"}
	topo, _ := New(original)

	copy1 := topo.Endpoints()
	copy1[0] = "mutated"

	copy2 := topo.Endpoints()
	if copy2[0] == "mutated" {
		t.Fatal("Endpoints should return an independent copy")
	}
}

func TestNext_ConcurrentAccess_NoPanic(t *testing.T) {
	topo, _ := New([]string{"h1:1", "h2:2", "h3:3"})

	const goroutines = 50
	const callsEach = 200

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < callsEach; j++ {
				_ = topo.Next()
			}
		}()
	}
	wg.Wait()
}
