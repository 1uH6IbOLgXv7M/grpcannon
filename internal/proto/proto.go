// Package proto provides utilities for resolving gRPC method descriptors
// from a target server using server reflection or a provided proto descriptor.
package proto

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// MethodDescriptor holds the resolved input/output type names for a gRPC method.
type MethodDescriptor struct {
	FullMethod  string
	ServiceName string
	MethodName  string
	InputType   string
	OutputType  string
}

// ParseFullMethod splits a full method string like "/pkg.Service/Method"
// into its service and method components.
func ParseFullMethod(fullMethod string) (service, method string, err error) {
	s := strings.TrimPrefix(fullMethod, "/")
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("proto: invalid method format %q, expected /Service/Method", fullMethod)
	}
	return parts[0], parts[1], nil
}

// ResolveMethod uses gRPC server reflection to look up the MethodDescriptor
// for the given fullMethod (e.g. "/helloworld.Greeter/SayHello").
func ResolveMethod(ctx context.Context, conn *grpc.ClientConn, fullMethod string) (*MethodDescriptor, error) {
	service, method, err := ParseFullMethod(fullMethod)
	if err != nil {
		return nil, err
	}

	stub := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
	stream, err := stub.ServerReflectionInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("proto: reflection stream: %w", err)
	}
	defer stream.CloseSend() //nolint:errcheck

	err = stream.Send(&grpc_reflection_v1alpha.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_FileContainingSymbol{
			FileContainingSymbol: service,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("proto: reflection send: %w", err)
	}

	resp, err := stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("proto: reflection recv: %w", err)
	}

	fdResp, ok := resp.MessageResponse.(*grpc_reflection_v1alpha.ServerReflectionResponse_FileDescriptorResponse)
	if !ok {
		return nil, fmt.Errorf("proto: unexpected reflection response type")
	}

	for _, b := range fdResp.FileDescriptorResponse.FileDescriptorProto {
		fd := &descriptorpb.FileDescriptorProto{}
		if err := proto.Unmarshal(b, fd); err != nil {
			continue
		}
		for _, svc := range fd.GetService() {
			if svc.GetName() != shortName(service) {
				continue
			}
			for _, m := range svc.GetMethod() {
				if m.GetName() == method {
					return &MethodDescriptor{
						FullMethod:  fullMethod,
						ServiceName: service,
						MethodName:  method,
						InputType:   m.GetInputType(),
						OutputType:  m.GetOutputType(),
					}, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("proto: method %q not found via reflection", fullMethod)
}

// shortName returns the last segment of a dot-separated name.
func shortName(s string) string {
	parts := strings.Split(s, ".")
	return parts[len(parts)-1]
}
