// Package label provides lightweight key-value metadata bags for annotating
// gRPC load test requests. Labels can be attached to individual requests or
// shared via a Registry so that concurrency profiles, warmup phases, and
// result snapshots can all carry consistent contextual metadata.
//
// Basic usage:
//
//	b := label.New("env", "prod", "region", "us-east-1")
//	v, ok := b.Get("env") // "prod", true
//
// Named sets can be stored and retrieved via a Registry:
//
//	reg := label.NewRegistry()
//	reg.Register("default", b)
//	bag, _ := reg.Lookup("default")
package label
