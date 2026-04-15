package invoker_test

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/yourorg/grpcannon/internal/invoker"
)

// dialLocalhost starts a bare gRPC server on a random port and returns a
// connected client connection together with a cleanup function.
func dialLocalhost(t *testing.T) (*grpc.ClientConn, func()) {
	t.Helper()
	srv := grpc.NewServer()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go srv.Serve(lis) //nolint:errcheck

	conn, err := grpc.NewClient(
		lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		srv.Stop()
		t.Fatalf("dial: %v", err)
	}
	return conn, func() {
		conn.Close()
		srv.Stop()
	}
}

func TestNew_NilConn_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil conn")
		}
	}()
	invoker.New(nil)
}

func TestDo_EmptyMethod_ReturnsError(t *testing.T) {
	conn, cleanup := dialLocalhost(t)
	defer cleanup()

	inv := invoker.New(conn)
	resp := inv.Do(context.Background(), invoker.Request{})
	if resp.Error == nil {
		t.Fatal("expected error for empty FullMethod")
	}
}

func TestDo_UnknownMethod_ReturnsError(t *testing.T) {
	conn, cleanup := dialLocalhost(t)
	defer cleanup()

	inv := invoker.New(conn)
	resp := inv.Do(context.Background(), invoker.Request{
		FullMethod: "/nonexistent.Service/Method",
	})
	if resp.Error == nil {
		t.Fatal("expected RPC error for unregistered method")
	}
	if resp.Duration <= 0 {
		t.Error("expected non-zero duration even on error")
	}
}

func TestDo_ContextCancelled_ReturnsError(t *testing.T) {
	conn, cleanup := dialLocalhost(t)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	inv := invoker.New(conn)
	resp := inv.Do(ctx, invoker.Request{FullMethod: "/pkg.Svc/M"})
	if resp.Error == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestDo_TimeoutApplied(t *testing.T) {
	conn, cleanup := dialLocalhost(t)
	defer cleanup()

	inv := invoker.New(conn)
	start := time.Now()
	resp := inv.Do(context.Background(), invoker.Request{
		FullMethod: "/pkg.Svc/Slow",
		Timeout:    50 * time.Millisecond,
	})
	elapsed := time.Since(start)

	if resp.Error == nil {
		t.Fatal("expected timeout error")
	}
	if elapsed > 500*time.Millisecond {
		t.Errorf("call took too long (%v); timeout not respected", elapsed)
	}
}
