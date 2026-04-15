package concurrency_test

import (
	"context"
	"testing"
	"time"

	"github.com/nicklaw5/grpcannon/internal/concurrency"
	"github.com/nicklaw5/grpcannon/internal/profile"
)

func stages(workers []int, dur time.Duration) []profile.Stage {
	out := make([]profile.Stage, len(workers))
	for i, w := range workers {
		out[i] = profile.Stage{Workers: w, Duration: dur}
	}
	return out
}

func TestController_EmitsAllStages(t *testing.T) {
	ctrl := concurrency.New(stages([]int{1, 2, 3}, 10*time.Millisecond))
	ctx := context.Background()

	var got []int
	go ctrl.Run(ctx)
	for w := range ctrl.Changes() {
		got = append(got, w)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(got))
	}
	for i, want := range []int{1, 2, 3} {
		if got[i] != want {
			t.Errorf("stage %d: want %d workers, got %d", i, want, got[i])
		}
	}
}

func TestController_ChannelClosedAfterRun(t *testing.T) {
	ctrl := concurrency.New(stages([]int{2}, 10*time.Millisecond))
	go ctrl.Run(context.Background())

	for range ctrl.Changes() {
	}
	// If the channel is not closed, this test would hang.
}

func TestController_CancelStopsEarly(t *testing.T) {
	ctrl := concurrency.New(stages([]int{1, 2, 3, 4, 5}, 200*time.Millisecond))
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var count int
	go ctrl.Run(ctx)
	for range ctrl.Changes() {
		count++
	}

	if count >= 5 {
		t.Errorf("expected cancellation to stop stages early, got %d", count)
	}
}

func TestController_NoStages_ClosesImmediately(t *testing.T) {
	ctrl := concurrency.New(nil)
	go ctrl.Run(context.Background())

	select {
	case _, ok := <-ctrl.Changes():
		if ok {
			t.Fatal("expected channel to be closed with no stages")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for channel close")
	}
}
