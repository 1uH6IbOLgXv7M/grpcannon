// Package reporter formats and prints load test results to an output stream.
package reporter

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/patrickward/grpcannon/internal/metrics"
)

// Reporter writes a formatted summary of load test results.
type Reporter struct {
	out io.Writer
}

// New creates a Reporter that writes to out.
func New(out io.Writer) *Reporter {
	return &Reporter{out: out}
}

// Print writes a human-readable latency histogram and summary to the output.
func (r *Reporter) Print(snap metrics.Snapshot, elapsed time.Duration) {
	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)

	fmt.Fprintf(r.out, "\n%s\n", strings.Repeat("─", 44))
	fmt.Fprintf(r.out, " gRPCannon Results\n")
	fmt.Fprintf(r.out, "%s\n\n", strings.Repeat("─", 44))

	fmt.Fprintf(w, "  Total requests:\t%d\n", snap.Total)
	fmt.Fprintf(w, "  Successful:\t%d\n", snap.Total-snap.Errors)
	fmt.Fprintf(w, "  Errors:\t%d\n", snap.Errors)
	fmt.Fprintf(w, "  Duration:\t%s\n", elapsed.Round(time.Millisecond))

	if snap.Total > 0 {
		rps := float64(snap.Total) / elapsed.Seconds()
		fmt.Fprintf(w, "  Req/sec:\t%.2f\n", rps)
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "  Latency (ms):\t\n")
	fmt.Fprintf(w, "    Mean:\t%.2f\n", float64(snap.Mean)/float64(time.Millisecond))
	fmt.Fprintf(w, "    P50:\t%.2f\n", float64(snap.P50)/float64(time.Millisecond))
	fmt.Fprintf(w, "    P90:\t%.2f\n", float64(snap.P90)/float64(time.Millisecond))
	fmt.Fprintf(w, "    P95:\t%.2f\n", float64(snap.P95)/float64(time.Millisecond))
	fmt.Fprintf(w, "    P99:\t%.2f\n", float64(snap.P99)/float64(time.Millisecond))

	w.Flush()
	fmt.Fprintf(r.out, "%s\n", strings.Repeat("─", 44))
}
