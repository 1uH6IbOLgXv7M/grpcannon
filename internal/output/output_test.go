package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/grpcannon/internal/metrics"
	"github.com/yourorg/grpcannon/internal/output"
)

func makeSnapshot() metrics.Snapshot {
	return metrics.Snapshot{
		Total:     100,
		Errors:    2,
		Duration:  5 * time.Second,
		ReqPerSec: 20.0,
		Mean:      50 * time.Millisecond,
		P50:       48 * time.Millisecond,
		P95:       95 * time.Millisecond,
		P99:       99 * time.Millisecond,
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	w := output.New(output.FormatText, nil)
	if w == nil {
		t.Fatal("expected non-nil Writer")
	}
}

func TestWrite_TextFormat_ContainsKeyFields(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.FormatText, &buf)

	if err := w.Write(makeSnapshot()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	for _, want := range []string{"Total: 100", "Errors: 2", "RPS: 20.00", "P50", "P95", "P99"} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q; got: %s", want, got)
		}
	}
}

func TestWrite_JSONFormat_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.FormatJSON, &buf)

	if err := w.Write(makeSnapshot()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}
}

func TestWrite_JSONFormat_ContainsTotalAndErrors(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.FormatJSON, &buf)

	snap := makeSnapshot()
	if err := w.Write(snap); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &result)

	if got, ok := result["Total"]; !ok || got == nil {
		t.Errorf("JSON output missing 'Total' field")
	}
	if got, ok := result["Errors"]; !ok || got == nil {
		t.Errorf("JSON output missing 'Errors' field")
	}
}
