// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: data_usage/data_usage.proto

package data_usage

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

// DataUsageSyncServiceClient is the client API for DataUsageSyncService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataUsageSyncServiceClient interface {
	CliVersionInformation(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*version.CliBuildInformation, error)
	SyncDataUsage(ctx context.Context, in *DataUsageSyncConfig, opts ...grpc.CallOption) (*DataUsageSyncResult, error)
}

type dataUsageSyncServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDataUsageSyncServiceClient(cc grpc.ClientConnInterface) DataUsageSyncServiceClient {
	return &dataUsageSyncServiceClient{cc}
}

func (c *dataUsageSyncServiceClient) CliVersionInformation(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*version.CliBuildInformation, error) {
	out := new(version.CliBuildInformation)
	err := c.cc.Invoke(ctx, "/data_usage.DataUsageSyncService/CliVersionInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataUsageSyncServiceClient) SyncDataUsage(ctx context.Context, in *DataUsageSyncConfig, opts ...grpc.CallOption) (*DataUsageSyncResult, error) {
	out := new(DataUsageSyncResult)
	err := c.cc.Invoke(ctx, "/data_usage.DataUsageSyncService/SyncDataUsage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataUsageSyncServiceServer is the server API for DataUsageSyncService service.
// All implementations must embed UnimplementedDataUsageSyncServiceServer
// for forward compatibility
type DataUsageSyncServiceServer interface {
	CliVersionInformation(context.Context, *emptypb.Empty) (*version.CliBuildInformation, error)
	SyncDataUsage(context.Context, *DataUsageSyncConfig) (*DataUsageSyncResult, error)
	mustEmbedUnimplementedDataUsageSyncServiceServer()
}

// UnimplementedDataUsageSyncServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDataUsageSyncServiceServer struct {
}

func (UnimplementedDataUsageSyncServiceServer) CliVersionInformation(context.Context, *emptypb.Empty) (*version.CliBuildInformation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CliVersionInformation not implemented")
}
func (UnimplementedDataUsageSyncServiceServer) SyncDataUsage(context.Context, *DataUsageSyncConfig) (*DataUsageSyncResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncDataUsage not implemented")
}
func (UnimplementedDataUsageSyncServiceServer) mustEmbedUnimplementedDataUsageSyncServiceServer() {}

// UnsafeDataUsageSyncServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataUsageSyncServiceServer will
// result in compilation errors.
type UnsafeDataUsageSyncServiceServer interface {
	mustEmbedUnimplementedDataUsageSyncServiceServer()
}

func RegisterDataUsageSyncServiceServer(s grpc.ServiceRegistrar, srv DataUsageSyncServiceServer) {
	s.RegisterService(&DataUsageSyncService_ServiceDesc, srv)
}

func _DataUsageSyncService_CliVersionInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataUsageSyncServiceServer).CliVersionInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/data_usage.DataUsageSyncService/CliVersionInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataUsageSyncServiceServer).CliVersionInformation(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataUsageSyncService_SyncDataUsage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DataUsageSyncConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataUsageSyncServiceServer).SyncDataUsage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/data_usage.DataUsageSyncService/SyncDataUsage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataUsageSyncServiceServer).SyncDataUsage(ctx, req.(*DataUsageSyncConfig))
	}
	return interceptor(ctx, in, info, handler)
}

// DataUsageSyncService_ServiceDesc is the grpc.ServiceDesc for DataUsageSyncService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataUsageSyncService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "data_usage.DataUsageSyncService",
	HandlerType: (*DataUsageSyncServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CliVersionInformation",
			Handler:    _DataUsageSyncService_CliVersionInformation_Handler,
		},
		{
			MethodName: "SyncDataUsage",
			Handler:    _DataUsageSyncService_SyncDataUsage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "data_usage/data_usage.proto",
}
