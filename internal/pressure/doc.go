// Package pressure derives a normalised back-pressure score [0.0, 1.0]
// from a sliding window of RPC latency and error observations.
//
// The composite score is computed as:
//
//	score = 0.7 * latencyScore + 0.3 * errorRate
//
// where latencyScore = min(meanLatency / HighLatency, 1.0).
//
// Observations older than the configured Window are automatically
// discarded on every call to Score or Record.
//
// Example:
//
//	tr := pressure.New(pressure.Config{
//		Window:      10 * time.Second,
//		HighLatency: 2 * time.Second,
//	})
//	tr.Record(350*time.Millisecond, false)
//	fmt.Println(tr.Score()) // e.g. 0.12
package pressure
