// Package proto provides helpers for resolving gRPC method descriptors
// via server reflection.
//
// It exposes ParseFullMethod for splitting a fully-qualified gRPC method
// string (e.g. "/helloworld.Greeter/SayHello") into its service and method
// components, and ResolveMethod for querying a live server's reflection API
// to obtain input/output type information for a given method.
//
// These utilities are used by the invoker package to construct dynamic
// gRPC requests without requiring pre-compiled proto stubs at build time.
package proto
