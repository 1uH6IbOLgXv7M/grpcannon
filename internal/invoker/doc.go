// Package invoker wraps a gRPC client connection and provides a simple
// Do method for executing unary calls with optional per-call timeouts and
// metadata headers.
//
// It uses a pass-through codec so callers can supply raw protobuf-encoded
// bytes without depending on generated code at the invoker layer. This
// keeps grpcannon codec-agnostic and suitable for black-box load testing.
//
// Typical usage:
//
//	inv := invoker.New(conn)
//	resp := inv.Do(ctx, invoker.Request{
//		FullMethod: "/mypackage.MyService/MyMethod",
//		Payload:    protoBytes,
//		Timeout:    2 * time.Second,
//	})
//	if resp.Error != nil {
//		// handle error
//	}
package invoker
