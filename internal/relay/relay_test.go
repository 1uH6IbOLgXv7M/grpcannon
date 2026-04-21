package relay_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/relay"
)

// makeSnapshot returns a minimal metrics snapshot value for use in tests.
func makeSnapshot(total, errors int64) relay.Snapshot {
	return relay.Snapshot{
		Total:  total,
		Errors: errors,
	}
}

func TestNew_NilInput_Panics(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil input channel")
		}
	}()
	relay.New(nil)
}

func TestRun_ForwardsAllSnapshots(t *testing.T) {
	t.Parallel()

	in := make(chan relay.Snapshot, 4)
	r := relay.New(in)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	out := r.Subscribe()

	snaps := []relay.Snapshot{
		makeSnapshot(1, 0),
		makeSnapshot(2, 1),
		makeSnapshot(3, 0),
	}
	for _, s := range snaps {
		in <- s
	}
	close(in)

	r.Run(ctx)

	var received []relay.Snapshot
	for s := range out {
		received = append(received, s)
	}

	if len(received) != len(snaps) {
		t.Fatalf("expected %d snapshots, got %d", len(snaps), len(received))
	}
	for i, s := range received {
		if s.Total != snaps[i].Total || s.Errors != snaps[i].Errors {
			t.Errorf("snapshot %d mismatch: got %+v, want %+v", i, s, snaps[i])
		}
	}
}

func TestRun_ContextCancel_StopsRelay(t *testing.T) {
	t.Parallel()

	in := make(chan relay.Snapshot) // unbuffered, nothing sent
	r := relay.New(in)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		defer close(done)
		r.Run(ctx)
	}()

	cancel()

	select {
	case <-done:
		// ok
	case <-time.After(time.Second):
		t.Fatal("Run did not stop after context cancellation")
	}
}

func TestSubscribe_MultipleSubscribers_AllReceive(t *testing.T) {
	t.Parallel()

	in := make(chan relay.Snapshot, 1)
	r := relay.New(in)

	out1 := r.Subscribe()
	out2 := r.Subscribe()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	snap := makeSnapshot(42, 3)
	in <- snap
	close(in)

	r.Run(ctx)

	for i, ch := range []<-chan relay.Snapshot{out1, out2} {
		var got relay.Snapshot
		select {
		case got = <-ch:
		default:
			t.Fatalf("subscriber %d received no snapshot", i)
		}
		if got.Total != snap.Total {
			t.Errorf("subscriber %d: got Total=%d, want %d", i, got.Total, snap.Total)
		}
	}
}

func TestRun_ConcurrentSubscribers_NoPanic(t *testing.T) {
	t.Parallel()

	in := make(chan relay.Snapshot, 8)
	r := relay.New(in)

	const numSubs = 10
	var wg sync.WaitGroup
	for i := 0; i < numSubs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.Subscribe()
		}()
	}
	wg.Wait()

	for i := 0; i < 8; i++ {
		in <- makeSnapshot(int64(i), 0)
	}
	close(in)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Should not panic.
	r.Run(ctx)
}
