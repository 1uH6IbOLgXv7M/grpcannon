package progress_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/bojand/grpcannon/internal/metrics"
	"github.com/bojand/grpcannon/internal/progress"
	"github.com/bojand/grpcannon/internal/snapshot"
)

func newCollector() *snapshot.Collector {
	rec := metrics.NewRecorder()
	return snapshot.NewCollector(rec)
}

func TestNew_DefaultsToStderr(t *testing.T) {
	c := newCollector()
	r := progress.New(nil, c, time.Second)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestNew_NegativeInterval_Defaults(t *testing.T) {
	c := newCollector()
	r := progress.New(nil, c, -1)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestStart_Stop_NoPanic(t *testing.T) {
	c := newCollector()
	var buf bytes.Buffer
	r := progress.New(&buf, c, 50*time.Millisecond)
	r.Start()
	time.Sleep(120 * time.Millisecond)
	r.Stop()
}

func TestPrint_ContainsExpectedFields(t *testing.T) {
	c := newCollector()
	var buf bytes.Buffer
	r := progress.New(&buf, c, 30*time.Millisecond)
	r.Start()
	time.Sleep(80 * time.Millisecond)
	r.Stop()

	out := buf.String()
	for _, want := range []string{"total=", "errors=", "err_rate=", "p50=", "p99="} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q; got: %s", want, out)
		}
	}
}

func TestStop_CalledTwice_Panics(t *testing.T) {
	// Stopping twice should panic because we close the channel; document that
	// callers must call Stop exactly once.
	c := newCollector()
	var buf bytes.Buffer
	r := progress.New(&buf, c, time.Second)
	r.Start()
	r.Stop()
	defer func() {
		if rec := recover(); rec == nil {
			t.Error("expected panic on double Stop")
		}
	}()
	r.Stop()
}
