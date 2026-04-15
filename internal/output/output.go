// Package output handles writing load test results to various destinations.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/yourorg/grpcannon/internal/metrics"
)

// Format represents the output format for results.
type Format string

const (
	// FormatText writes a human-readable summary.
	FormatText Format = "text"
	// FormatJSON writes a machine-readable JSON object.
	FormatJSON Format = "json"
)

// Writer writes a metrics snapshot to an output destination.
type Writer struct {
	format Format
	w      io.Writer
}

// New creates a Writer that writes to w using the given format.
// If w is nil, os.Stdout is used.
func New(format Format, w io.Writer) *Writer {
	if w == nil {
		w = os.Stdout
	}
	return &Writer{format: format, w: w}
}

// Write serialises the snapshot according to the configured format.
func (wr *Writer) Write(snap metrics.Snapshot) error {
	switch wr.format {
	case FormatJSON:
		return wr.writeJSON(snap)
	default:
		return wr.writeText(snap)
	}
}

func (wr *Writer) writeText(snap metrics.Snapshot) error {
	_, err := fmt.Fprintf(wr.w,
		"Total: %d  Errors: %d  Duration: %s  RPS: %.2f  P50: %s  P95: %s  P99: %s\n",
		snap.Total,
		snap.Errors,
		snap.Duration,
		snap.ReqPerSec,
		snap.P50,
		snap.P95,
		snap.P99,
	)
	return err
}

func (wr *Writer) writeJSON(snap metrics.Snapshot) error {
	enc := json.NewEncoder(wr.w)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}
