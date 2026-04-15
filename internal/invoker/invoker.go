// Package invoker provides a generic gRPC method invoker that executes
// unary calls against a target service using a pre-established connection.
package invoker

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Request holds the parameters for a single gRPC invocation.
type Request struct {
	// FullMethod is the fully-qualified gRPC method name, e.g. "/pkg.Service/Method".
	FullMethod string
	// Payload is the raw protobuf-encoded request message.
	Payload []byte
	// Metadata contains optional key-value pairs sent as gRPC headers.
	Metadata map[string]string
	// Timeout caps the individual call duration. Zero means no per-call timeout.
	Timeout time.Duration
}

// Response captures the outcome of a single gRPC invocation.
type Response struct {
	// Duration is the wall waiting for the call to complete.
	Duration time.Duration
	// Error is non-nil when the call failed.
	n
// Invoker executes gRPC unary calls over a shared connection.
type Invoker struct {
	conn *grpc.ClientConn
}

// New creates an Invoker that reuses conn for all calls.
func New(conn *grpc.ClientConn) *Invoker {
	if conn == nil {
		panic("invoker: conn must not be nil")
	}
	return &Invoker{conn: conn}
}

// Do executes a single unary gRPC call described by req.
// The parent context controls overall cancellation; req.Timeout adds an
// additional deadline scoped to this call only.
func (inv *Invoker) Do(ctx context.Context, req Request) Response {
	if req.FullMethod == "" {
		return Response{Error: fmt.Errorf("invoker: FullMethod must not be empty")}
	}

	callCtx := ctx
	var cancel context.CancelFunc
	if req.Timeout > 0 {
		callCtx, cancel = context.WithTimeout(ctx, req.Timeout)
		defer cancel()
	}

	if len(req.Metadata) > 0 {
		md := metadata.New(req.Metadata)
		callCtx = metadata.NewOutgoingContext(callCtx, md)
	}

	var reply []byte
	start := time.Now()
	err := inv.conn.Invoke(callCtx, req.FullMethod, req.Payload, &reply, grpc.ForceCodec(rawCodec{}))
	elapsed := time.Since(start)

	return Response{Duration: elapsed, Error: err}
}
