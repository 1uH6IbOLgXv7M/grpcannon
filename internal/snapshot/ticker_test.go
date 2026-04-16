package snapshot_test

import (
	"context"
	"testing"
	"time"

	"github.com/nicklaw5/grpcannon/internal/metrics"
	"github.com/nicklaw5/grpcannon/internal/snapshot"
)

func TestTicker_DeliversSnapshots(t *testing.T) {
	r := metrics.NewRecorder()
	c := snapshot.NewCollector(r)
	tk := snapshot.NewTicker(c, 10*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Millisecond)
	defer cancel()

	go tk.Run(ctx)

	var count int
	for range tk.C() {
		count++
	}

	if count < 3 {
		t.Fatalf("expected at least 3 snapshots, got %d", count)
	}
}

func TestTicker_ChannelClosedAfterCancel(t *testing.T) {
	r := metrics.NewRecorder()
	c := snapshot.NewCollector(r)
	tk := snapshot.NewTicker(c, 10*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	go tk.Run(ctx)
	time.Sleep(15 * time.Millisecond)
	cancel()

	// drain until closed
	timeout := time.After(200 * time.Millisecond)
	for {
		select {
		case _, ok := <-tk.C():
			if !ok {
				return
			}
		case <-timeout:
			t.Fatal("channel not closed after cancel")
		}
	}
}
