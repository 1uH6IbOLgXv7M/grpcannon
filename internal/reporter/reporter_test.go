package reporter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/patrickward/grpcannon/internal/metrics"
	"github.com/patrickward/grpcannon/internal/reporter"
)

func snapshot() metrics.Snapshot {
	return metrics.Snapshot{
		Total:  100,
		Errors: 3,
		Mean:   12 * time.Millisecond,
		P50:    10 * time.Millisecond,
		P90:    20 * time.Millisecond,
		P95:    25 * time.Millisecond,
		P99:    40 * time.Millisecond,
	}
}

func TestPrint_ContainsTotals(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)
	r.Print(snapshot(), 5*time.Second)

	out := buf.String()
	for _, want := range []string{"100", "3", "97"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q\ngot:\n%s", want, out)
		}
	}
}

func TestPrint_ContainsLatencyPercentiles(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)
	r.Print(snapshot(), 5*time.Second)

	out := buf.String()
	for _, want := range []string{"P50", "P90", "P95", "P99", "Mean"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q\ngot:\n%s", want, out)
		}
	}
}

func TestPrint_ContainsReqPerSec(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)
	r.Print(snapshot(), 10*time.Second)

	out := buf.String()
	// 100 requests / 10s = 10.00 req/sec
	if !strings.Contains(out, "10.00") {
		t.Errorf("expected req/sec of 10.00 in output\ngot:\n%s", out)
	}
}

func TestPrint_ZeroTotal_NoReqSec(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)
	r.Print(metrics.Snapshot{}, 5*time.Second)

	out := buf.String()
	if strings.Contains(out, "Req/sec") {
		t.Errorf("expected no Req/sec line for zero requests\ngot:\n%s", out)
	}
}
