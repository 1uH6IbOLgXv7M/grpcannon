// Package shedder provides adaptive load shedding for gRPC request pipelines.
//
// A Shedder combines two complementary strategies:
//
//  1. In-flight ceiling – hard cap on the number of concurrent requests.
//     When the ceiling is reached every new request is immediately rejected
//     with ErrShed so the backend is never overwhelmed.
//
//  2. Error-rate gate – a rolling-window error rate check.  Once the observed
//     error rate climbs above the configured threshold new requests are shed
//     until the rate recovers, giving the backend time to stabilise.
//
// Both strategies can be used independently or together.
package shedder
