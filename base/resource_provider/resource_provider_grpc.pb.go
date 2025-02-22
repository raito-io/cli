// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: resource_provider/resource_provider.proto

package resource_provider

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

// ResourceProviderServiceClient is the client API for ResourceProviderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ResourceProviderServiceClient interface {
	CliVersionInformation(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*version.CliBuildInformation, error)
	UpdateResources(ctx context.Context, in *UpdateResourceInput, opts ...grpc.CallOption) (*UpdateResourceResult, error)
}

type resourceProviderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewResourceProviderServiceClient(cc grpc.ClientConnInterface) ResourceProviderServiceClient {
	return &resourceProviderServiceClient{cc}
}

func (c *resourceProviderServiceClient) CliVersionInformation(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*version.CliBuildInformation, error) {
	out := new(version.CliBuildInformation)
	err := c.cc.Invoke(ctx, "/resource_provider.ResourceProviderService/CliVersionInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceProviderServiceClient) UpdateResources(ctx context.Context, in *UpdateResourceInput, opts ...grpc.CallOption) (*UpdateResourceResult, error) {
	out := new(UpdateResourceResult)
	err := c.cc.Invoke(ctx, "/resource_provider.ResourceProviderService/UpdateResources", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceProviderServiceServer is the server API for ResourceProviderService service.
// All implementations must embed UnimplementedResourceProviderServiceServer
// for forward compatibility
type ResourceProviderServiceServer interface {
	CliVersionInformation(context.Context, *emptypb.Empty) (*version.CliBuildInformation, error)
	UpdateResources(context.Context, *UpdateResourceInput) (*UpdateResourceResult, error)
	mustEmbedUnimplementedResourceProviderServiceServer()
}

// UnimplementedResourceProviderServiceServer must be embedded to have forward compatible implementations.
type UnimplementedResourceProviderServiceServer struct {
}

func (UnimplementedResourceProviderServiceServer) CliVersionInformation(context.Context, *emptypb.Empty) (*version.CliBuildInformation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CliVersionInformation not implemented")
}
func (UnimplementedResourceProviderServiceServer) UpdateResources(context.Context, *UpdateResourceInput) (*UpdateResourceResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateResources not implemented")
}
func (UnimplementedResourceProviderServiceServer) mustEmbedUnimplementedResourceProviderServiceServer() {
}

// UnsafeResourceProviderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ResourceProviderServiceServer will
// result in compilation errors.
type UnsafeResourceProviderServiceServer interface {
	mustEmbedUnimplementedResourceProviderServiceServer()
}

func RegisterResourceProviderServiceServer(s grpc.ServiceRegistrar, srv ResourceProviderServiceServer) {
	s.RegisterService(&ResourceProviderService_ServiceDesc, srv)
}

func _ResourceProviderService_CliVersionInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceProviderServiceServer).CliVersionInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/resource_provider.ResourceProviderService/CliVersionInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceProviderServiceServer).CliVersionInformation(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ResourceProviderService_UpdateResources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateResourceInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceProviderServiceServer).UpdateResources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/resource_provider.ResourceProviderService/UpdateResources",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceProviderServiceServer).UpdateResources(ctx, req.(*UpdateResourceInput))
	}
	return interceptor(ctx, in, info, handler)
}

// ResourceProviderService_ServiceDesc is the grpc.ServiceDesc for ResourceProviderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ResourceProviderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "resource_provider.ResourceProviderService",
	HandlerType: (*ResourceProviderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CliVersionInformation",
			Handler:    _ResourceProviderService_CliVersionInformation_Handler,
		},
		{
			MethodName: "UpdateResources",
			Handler:    _ResourceProviderService_UpdateResources_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "resource_provider/resource_provider.proto",
}
