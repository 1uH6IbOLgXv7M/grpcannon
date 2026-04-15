package dialer_test

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/yourorg/grpcannon/internal/dialer"
)

// startEchoServer spins up a bare gRPC server on a random local port and
// returns its address together with a stop function.
func startEchoServer(t *testing.T) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	srv := grpc.NewServer()
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)
	return lis.Addr().String()
}

func TestConnect_EmptyTarget_ReturnsError(t *testing.T) {
	_, err := dialer.Connect(context.Background(), "", dialer.Options{Insecure: true})
	if err == nil {
		t.Fatal("expected error for empty target, got nil")
	}
}

func TestConnect_Insecure_Success(t *testing.T) {
	addr := startEchoServer(t)

	conn, err := dialer.Connect(context.Background(), addr, dialer.Options{
		Insecure: true,
		Timeout:  2 * time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer conn.Close()

	if conn == nil {
		t.Fatal("expected non-nil connection")
	}
}

func TestConnect_Timeout_ReturnsError(t *testing.T) {
	// Nothing listening on this port — dial should time out.
	_, err := dialer.Connect(context.Background(), "127.0.0.1:1", dialer.Options{
		Insecure: true,
		Timeout:  100 * time.Millisecond,
	})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestConnect_CancelledContext_ReturnsError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err := dialer.Connect(ctx, "127.0.0.1:1", dialer.Options{Insecure: true})
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
}

// Ensure the package-level grpc import is used (compile check).
var _ = grpc.WithTransportCredentials(insecure.NewCredentials())
