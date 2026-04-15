// Package sampler implements reservoir sampling (Algorithm R) for collecting
// a bounded, statistically representative set of gRPC response payloads
// observed during a grpcannon load test run.
//
// Usage:
//
//	s := sampler.New(100, time.Now().UnixNano())
//	s.Add("/pkg.Service/Method", responseBytes)
//	for _, sample := range s.Samples() {
//		fmt.Println(sample.Method, sample.Payload)
//	}
//
// The sampler is safe for concurrent use by multiple goroutines.
package sampler
