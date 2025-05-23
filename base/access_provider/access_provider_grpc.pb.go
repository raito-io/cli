// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: access_provider/access_provider.proto

package access_provider

import (
	context "context"
	version "github.com/raito-io/cli/base/util/version"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AccessProviderSyncServiceClient is the client API for AccessProviderSyncService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccessProviderSyncServiceClient interface {
	CliVersionInformation(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*version.CliBuildInformation, error)
	SyncFromTarget(ctx context.Context, in *AccessSyncFromTarget, opts ...grpc.CallOption) (*AccessSyncResult, error)
	SyncToTarget(ctx context.Context, in *AccessSyncToTarget, opts ...grpc.CallOption) (*AccessSyncResult, error)
	SyncConfig(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AccessSyncConfig, error)
}

type accessProviderSyncServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccessProviderSyncServiceClient(cc grpc.ClientConnInterface) AccessProviderSyncServiceClient {
	return &accessProviderSyncServiceClient{cc}
}

func (c *accessProviderSyncServiceClient) CliVersionInformation(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*version.CliBuildInformation, error) {
	out := new(version.CliBuildInformation)
	err := c.cc.Invoke(ctx, "/access_provider.AccessProviderSyncService/CliVersionInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessProviderSyncServiceClient) SyncFromTarget(ctx context.Context, in *AccessSyncFromTarget, opts ...grpc.CallOption) (*AccessSyncResult, error) {
	out := new(AccessSyncResult)
	err := c.cc.Invoke(ctx, "/access_provider.AccessProviderSyncService/SyncFromTarget", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessProviderSyncServiceClient) SyncToTarget(ctx context.Context, in *AccessSyncToTarget, opts ...grpc.CallOption) (*AccessSyncResult, error) {
	out := new(AccessSyncResult)
	err := c.cc.Invoke(ctx, "/access_provider.AccessProviderSyncService/SyncToTarget", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessProviderSyncServiceClient) SyncConfig(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AccessSyncConfig, error) {
	out := new(AccessSyncConfig)
	err := c.cc.Invoke(ctx, "/access_provider.AccessProviderSyncService/SyncConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccessProviderSyncServiceServer is the server API for AccessProviderSyncService service.
// All implementations must embed UnimplementedAccessProviderSyncServiceServer
// for forward compatibility
type AccessProviderSyncServiceServer interface {
	CliVersionInformation(context.Context, *emptypb.Empty) (*version.CliBuildInformation, error)
	SyncFromTarget(context.Context, *AccessSyncFromTarget) (*AccessSyncResult, error)
	SyncToTarget(context.Context, *AccessSyncToTarget) (*AccessSyncResult, error)
	SyncConfig(context.Context, *emptypb.Empty) (*AccessSyncConfig, error)
	mustEmbedUnimplementedAccessProviderSyncServiceServer()
}

// UnimplementedAccessProviderSyncServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAccessProviderSyncServiceServer struct {
}

func (UnimplementedAccessProviderSyncServiceServer) CliVersionInformation(context.Context, *emptypb.Empty) (*version.CliBuildInformation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CliVersionInformation not implemented")
}
func (UnimplementedAccessProviderSyncServiceServer) SyncFromTarget(context.Context, *AccessSyncFromTarget) (*AccessSyncResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncFromTarget not implemented")
}
func (UnimplementedAccessProviderSyncServiceServer) SyncToTarget(context.Context, *AccessSyncToTarget) (*AccessSyncResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncToTarget not implemented")
}
func (UnimplementedAccessProviderSyncServiceServer) SyncConfig(context.Context, *emptypb.Empty) (*AccessSyncConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncConfig not implemented")
}
func (UnimplementedAccessProviderSyncServiceServer) mustEmbedUnimplementedAccessProviderSyncServiceServer() {
}

// UnsafeAccessProviderSyncServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccessProviderSyncServiceServer will
// result in compilation errors.
type UnsafeAccessProviderSyncServiceServer interface {
	mustEmbedUnimplementedAccessProviderSyncServiceServer()
}

func RegisterAccessProviderSyncServiceServer(s grpc.ServiceRegistrar, srv AccessProviderSyncServiceServer) {
	s.RegisterService(&AccessProviderSyncService_ServiceDesc, srv)
}

func _AccessProviderSyncService_CliVersionInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessProviderSyncServiceServer).CliVersionInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/access_provider.AccessProviderSyncService/CliVersionInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessProviderSyncServiceServer).CliVersionInformation(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessProviderSyncService_SyncFromTarget_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccessSyncFromTarget)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessProviderSyncServiceServer).SyncFromTarget(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/access_provider.AccessProviderSyncService/SyncFromTarget",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessProviderSyncServiceServer).SyncFromTarget(ctx, req.(*AccessSyncFromTarget))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessProviderSyncService_SyncToTarget_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccessSyncToTarget)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessProviderSyncServiceServer).SyncToTarget(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/access_provider.AccessProviderSyncService/SyncToTarget",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessProviderSyncServiceServer).SyncToTarget(ctx, req.(*AccessSyncToTarget))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessProviderSyncService_SyncConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessProviderSyncServiceServer).SyncConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/access_provider.AccessProviderSyncService/SyncConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessProviderSyncServiceServer).SyncConfig(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// AccessProviderSyncService_ServiceDesc is the grpc.ServiceDesc for AccessProviderSyncService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccessProviderSyncService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "access_provider.AccessProviderSyncService",
	HandlerType: (*AccessProviderSyncServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CliVersionInformation",
			Handler:    _AccessProviderSyncService_CliVersionInformation_Handler,
		},
		{
			MethodName: "SyncFromTarget",
			Handler:    _AccessProviderSyncService_SyncFromTarget_Handler,
		},
		{
			MethodName: "SyncToTarget",
			Handler:    _AccessProviderSyncService_SyncToTarget_Handler,
		},
		{
			MethodName: "SyncConfig",
			Handler:    _AccessProviderSyncService_SyncConfig_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "access_provider/access_provider.proto",
}
