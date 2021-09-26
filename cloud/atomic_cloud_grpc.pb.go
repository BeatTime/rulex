// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package cloud

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AtomicCloudServiceClient is the client API for AtomicCloudService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AtomicCloudServiceClient interface {
	CallCloud(ctx context.Context, in *Service, opts ...grpc.CallOption) (*CallResult, error)
}

type atomicCloudServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAtomicCloudServiceClient(cc grpc.ClientConnInterface) AtomicCloudServiceClient {
	return &atomicCloudServiceClient{cc}
}

func (c *atomicCloudServiceClient) CallCloud(ctx context.Context, in *Service, opts ...grpc.CallOption) (*CallResult, error) {
	out := new(CallResult)
	err := c.cc.Invoke(ctx, "/cloud.AtomicCloudService/CallCloud", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AtomicCloudServiceServer is the server API for AtomicCloudService service.
// All implementations must embed UnimplementedAtomicCloudServiceServer
// for forward compatibility
type AtomicCloudServiceServer interface {
	CallCloud(context.Context, *Service) (*CallResult, error)
	mustEmbedUnimplementedAtomicCloudServiceServer()
}

// UnimplementedAtomicCloudServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAtomicCloudServiceServer struct {
}

func (UnimplementedAtomicCloudServiceServer) CallCloud(context.Context, *Service) (*CallResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CallCloud not implemented")
}
func (UnimplementedAtomicCloudServiceServer) mustEmbedUnimplementedAtomicCloudServiceServer() {}

// UnsafeAtomicCloudServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AtomicCloudServiceServer will
// result in compilation errors.
type UnsafeAtomicCloudServiceServer interface {
	mustEmbedUnimplementedAtomicCloudServiceServer()
}

func RegisterAtomicCloudServiceServer(s grpc.ServiceRegistrar, srv AtomicCloudServiceServer) {
	s.RegisterService(&AtomicCloudService_ServiceDesc, srv)
}

func _AtomicCloudService_CallCloud_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Service)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AtomicCloudServiceServer).CallCloud(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloud.AtomicCloudService/CallCloud",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AtomicCloudServiceServer).CallCloud(ctx, req.(*Service))
	}
	return interceptor(ctx, in, info, handler)
}

// AtomicCloudService_ServiceDesc is the grpc.ServiceDesc for AtomicCloudService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AtomicCloudService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cloud.AtomicCloudService",
	HandlerType: (*AtomicCloudServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CallCloud",
			Handler:    _AtomicCloudService_CallCloud_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "atomic_cloud.proto",
}