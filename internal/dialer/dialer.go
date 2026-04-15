// Package dialer provides gRPC connection helpers for grpcannon.
package dialer

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Options controls how a gRPC connection is established.
type Options struct {
	// Insecure disables TLS. Mutually exclusive with TLSConfig.
	Insecure bool

	// TLSConfig is an optional custom TLS configuration.
	TLSConfig *tls.Config

	// Timeout is the maximum time to wait for the connection to be ready.
	// Zero means no timeout.
	Timeout time.Duration
}

// Connect dials the given target and returns a ready gRPC ClientConn.
// The caller is responsible for closing the connection.
func Connect(ctx context.Context, target string, opts Options) (*grpc.ClientConn, error) {
	if target == "" {
		return nil, fmt.Errorf("dialer: target must not be empty")
	}

	var creds credentials.TransportCredentials
	switch {
	case opts.Insecure:
		creds = insecure.NewCredentials()
	case opts.TLSConfig != nil:
		creds = credentials.NewTLS(opts.TLSConfig)
	default:
		creds = credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})
	}

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
	}

	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	conn, err := grpc.DialContext(ctx, target, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("dialer: connect to %q: %w", target, err)
	}

	return conn, nil
}
