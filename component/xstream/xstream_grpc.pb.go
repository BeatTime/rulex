// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: xstream.proto

package xstream

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

const (
	XStream_OnApproached_FullMethodName = "/xstream.XStream/OnApproached"
	XStream_SendStream_FullMethodName   = "/xstream.XStream/SendStream"
)

// XStreamClient is the client API for XStream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type XStreamClient interface {
	// 收到来自其他端点的请求
	OnApproached(ctx context.Context, opts ...grpc.CallOption) (XStream_OnApproachedClient, error)
	// 给其他端点发送请求
	SendStream(ctx context.Context, in *Request, opts ...grpc.CallOption) (XStream_SendStreamClient, error)
}

type xStreamClient struct {
	cc grpc.ClientConnInterface
}

func NewXStreamClient(cc grpc.ClientConnInterface) XStreamClient {
	return &xStreamClient{cc}
}

func (c *xStreamClient) OnApproached(ctx context.Context, opts ...grpc.CallOption) (XStream_OnApproachedClient, error) {
	stream, err := c.cc.NewStream(ctx, &XStream_ServiceDesc.Streams[0], XStream_OnApproached_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &xStreamOnApproachedClient{stream}
	return x, nil
}

type XStream_OnApproachedClient interface {
	Send(*Request) error
	CloseAndRecv() (*Request, error)
	grpc.ClientStream
}

type xStreamOnApproachedClient struct {
	grpc.ClientStream
}

func (x *xStreamOnApproachedClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *xStreamOnApproachedClient) CloseAndRecv() (*Request, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Request)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *xStreamClient) SendStream(ctx context.Context, in *Request, opts ...grpc.CallOption) (XStream_SendStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &XStream_ServiceDesc.Streams[1], XStream_SendStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &xStreamSendStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type XStream_SendStreamClient interface {
	Recv() (*Response, error)
	grpc.ClientStream
}

type xStreamSendStreamClient struct {
	grpc.ClientStream
}

func (x *xStreamSendStreamClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// XStreamServer is the server API for XStream service.
// All implementations must embed UnimplementedXStreamServer
// for forward compatibility
type XStreamServer interface {
	// 收到来自其他端点的请求
	OnApproached(XStream_OnApproachedServer) error
	// 给其他端点发送请求
	SendStream(*Request, XStream_SendStreamServer) error
	mustEmbedUnimplementedXStreamServer()
}

// UnimplementedXStreamServer must be embedded to have forward compatible implementations.
type UnimplementedXStreamServer struct {
}

func (UnimplementedXStreamServer) OnApproached(XStream_OnApproachedServer) error {
	return status.Errorf(codes.Unimplemented, "method OnApproached not implemented")
}
func (UnimplementedXStreamServer) SendStream(*Request, XStream_SendStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method SendStream not implemented")
}
func (UnimplementedXStreamServer) mustEmbedUnimplementedXStreamServer() {}

// UnsafeXStreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to XStreamServer will
// result in compilation errors.
type UnsafeXStreamServer interface {
	mustEmbedUnimplementedXStreamServer()
}

func RegisterXStreamServer(s grpc.ServiceRegistrar, srv XStreamServer) {
	s.RegisterService(&XStream_ServiceDesc, srv)
}

func _XStream_OnApproached_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(XStreamServer).OnApproached(&xStreamOnApproachedServer{stream})
}

type XStream_OnApproachedServer interface {
	SendAndClose(*Request) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type xStreamOnApproachedServer struct {
	grpc.ServerStream
}

func (x *xStreamOnApproachedServer) SendAndClose(m *Request) error {
	return x.ServerStream.SendMsg(m)
}

func (x *xStreamOnApproachedServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _XStream_SendStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(XStreamServer).SendStream(m, &xStreamSendStreamServer{stream})
}

type XStream_SendStreamServer interface {
	Send(*Response) error
	grpc.ServerStream
}

type xStreamSendStreamServer struct {
	grpc.ServerStream
}

func (x *xStreamSendStreamServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

// XStream_ServiceDesc is the grpc.ServiceDesc for XStream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var XStream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "xstream.XStream",
	HandlerType: (*XStreamServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "OnApproached",
			Handler:       _XStream_OnApproached_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "SendStream",
			Handler:       _XStream_SendStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "xstream.proto",
}