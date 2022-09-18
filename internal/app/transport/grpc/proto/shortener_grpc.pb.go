// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: internal/app/transport/grpc/proto/shortener.proto

package pb

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

// ShortenerClient is the client API for Shortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortenerClient interface {
	Ping(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*PingResponse, error)
	HandlePost(ctx context.Context, in *HandlePostRequest, opts ...grpc.CallOption) (*HandlePostResponse, error)
	HandleGet(ctx context.Context, in *HandleGetRequest, opts ...grpc.CallOption) (*HandleGetResponse, error)
	HandleGetUserURLs(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*HandleGetUserURLsResponse, error)
	HandlePostShortenBatch(ctx context.Context, in *HandlePostShortenBatchRequest, opts ...grpc.CallOption) (*HandlePostShortenBatchResponse, error)
}

type shortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewShortenerClient(cc grpc.ClientConnInterface) ShortenerClient {
	return &shortenerClient{cc}
}

func (c *shortenerClient) Ping(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) HandlePost(ctx context.Context, in *HandlePostRequest, opts ...grpc.CallOption) (*HandlePostResponse, error) {
	out := new(HandlePostResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/HandlePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) HandleGet(ctx context.Context, in *HandleGetRequest, opts ...grpc.CallOption) (*HandleGetResponse, error) {
	out := new(HandleGetResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/HandleGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) HandleGetUserURLs(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*HandleGetUserURLsResponse, error) {
	out := new(HandleGetUserURLsResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/HandleGetUserURLs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) HandlePostShortenBatch(ctx context.Context, in *HandlePostShortenBatchRequest, opts ...grpc.CallOption) (*HandlePostShortenBatchResponse, error) {
	out := new(HandlePostShortenBatchResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/HandlePostShortenBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortenerServer is the server API for Shortener service.
// All implementations must embed UnimplementedShortenerServer
// for forward compatibility
type ShortenerServer interface {
	Ping(context.Context, *EmptyRequest) (*PingResponse, error)
	HandlePost(context.Context, *HandlePostRequest) (*HandlePostResponse, error)
	HandleGet(context.Context, *HandleGetRequest) (*HandleGetResponse, error)
	HandleGetUserURLs(context.Context, *EmptyRequest) (*HandleGetUserURLsResponse, error)
	HandlePostShortenBatch(context.Context, *HandlePostShortenBatchRequest) (*HandlePostShortenBatchResponse, error)
	mustEmbedUnimplementedShortenerServer()
}

// UnimplementedShortenerServer must be embedded to have forward compatible implementations.
type UnimplementedShortenerServer struct {
}

func (UnimplementedShortenerServer) Ping(context.Context, *EmptyRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedShortenerServer) HandlePost(context.Context, *HandlePostRequest) (*HandlePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandlePost not implemented")
}
func (UnimplementedShortenerServer) HandleGet(context.Context, *HandleGetRequest) (*HandleGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleGet not implemented")
}
func (UnimplementedShortenerServer) HandleGetUserURLs(context.Context, *EmptyRequest) (*HandleGetUserURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleGetUserURLs not implemented")
}
func (UnimplementedShortenerServer) HandlePostShortenBatch(context.Context, *HandlePostShortenBatchRequest) (*HandlePostShortenBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandlePostShortenBatch not implemented")
}
func (UnimplementedShortenerServer) mustEmbedUnimplementedShortenerServer() {}

// UnsafeShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortenerServer will
// result in compilation errors.
type UnsafeShortenerServer interface {
	mustEmbedUnimplementedShortenerServer()
}

func RegisterShortenerServer(s grpc.ServiceRegistrar, srv ShortenerServer) {
	s.RegisterService(&Shortener_ServiceDesc, srv)
}

func _Shortener_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).Ping(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_HandlePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HandlePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).HandlePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/HandlePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).HandlePost(ctx, req.(*HandlePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_HandleGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HandleGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).HandleGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/HandleGet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).HandleGet(ctx, req.(*HandleGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_HandleGetUserURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).HandleGetUserURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/HandleGetUserURLs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).HandleGetUserURLs(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_HandlePostShortenBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HandlePostShortenBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).HandlePostShortenBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/HandlePostShortenBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).HandlePostShortenBatch(ctx, req.(*HandlePostShortenBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Shortener_ServiceDesc is the grpc.ServiceDesc for Shortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.Shortener",
	HandlerType: (*ShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Shortener_Ping_Handler,
		},
		{
			MethodName: "HandlePost",
			Handler:    _Shortener_HandlePost_Handler,
		},
		{
			MethodName: "HandleGet",
			Handler:    _Shortener_HandleGet_Handler,
		},
		{
			MethodName: "HandleGetUserURLs",
			Handler:    _Shortener_HandleGetUserURLs_Handler,
		},
		{
			MethodName: "HandlePostShortenBatch",
			Handler:    _Shortener_HandlePostShortenBatch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/app/transport/grpc/proto/shortener.proto",
}
