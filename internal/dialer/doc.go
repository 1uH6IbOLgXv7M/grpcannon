// Package dialer abstracts gRPC connection establishment for grpcannon.
//
// It provides a single Connect function that accepts a target address and
// an Options struct, handling TLS negotiation, insecure mode, and dial
// timeouts in a consistent way across the rest of the application.
//
// Usage:
//
//	conn, err := dialer.Connect(ctx, "localhost:50051", dialer.Options{
//		Insecure: true,
//		Timeout:  5 * time.Second,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer conn.Close()
package dialer
