package tee_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/metrics"
	"github.com/example/grpcannon/internal/tee"
)

type countSink struct{ n atomic.Int64 }

func (c *countSink) Write(_ metrics.Snapshot) { c.n.Add(1) }

func TestSend_DeliveriesToAllSinks(t *testing.T) {
	a, b := &countSink{}, &countSink{}
	tr := tee.New(a, b)
	tr.Send(metrics.Snapshot{})
	if a.n.Load() != 1 || b.n.Load() != 1 {
		t.Fatalf("expected each sink to receive 1 snapshot, got %d %d", a.n.Load(), b.n.Load())
	}
}

func TestAdd_NewSinkReceivesSubsequentSnapshots(t *testing.T) {
	a := &countSink{}
	tr := tee.New(a)
	b := &countSink{}
	tr.Add(b)
	tr.Send(metrics.Snapshot{})
	if b.n.Load() != 1 {
		t.Fatalf("expected late-added sink to receive snapshot")
	}
}

func TestRun_DeliversThenExitsOnChannelClose(t *testing.T) {
	s := &countSink{}
	tr := tee.New(s)
	ch := make(chan metrics.Snapshot, 3)
	for i := 0; i < 3; i++ {
		ch <- metrics.Snapshot{}
	}
	close(ch)
	tr.Run(context.Background(), ch)
	if s.n.Load() != 3 {
		t.Fatalf("expected 3 deliveries, got %d", s.n.Load())
	}
}

func TestRun_ExitsOnContextCancel(t *testing.T) {
	s := &countSink{}
	tr := tee.New(s)
	ch := make(chan metrics.Snapshot)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		tr.Run(ctx, ch)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not exit after context cancellation")
	}
}

func TestSend_NoSinks_NoPanic(t *testing.T) {
	tr := tee.New()
	tr.Send(metrics.Snapshot{})
}
