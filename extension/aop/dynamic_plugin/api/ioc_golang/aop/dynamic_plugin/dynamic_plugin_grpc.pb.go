// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package dynamic_plugin

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

// DynamicPluginServiceClient is the client API for DynamicPluginService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DynamicPluginServiceClient interface {
	Update(ctx context.Context, in *DynamicPluginUpdateRequest, opts ...grpc.CallOption) (*DynamicPluginUpdateResponse, error)
}

type dynamicPluginServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDynamicPluginServiceClient(cc grpc.ClientConnInterface) DynamicPluginServiceClient {
	return &dynamicPluginServiceClient{cc}
}

func (c *dynamicPluginServiceClient) Update(ctx context.Context, in *DynamicPluginUpdateRequest, opts ...grpc.CallOption) (*DynamicPluginUpdateResponse, error) {
	out := new(DynamicPluginUpdateResponse)
	err := c.cc.Invoke(ctx, "/ioc_golang.aop.dynamic_plugin.DynamicPluginService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DynamicPluginServiceServer is the server API for DynamicPluginService service.
// All implementations must embed UnimplementedDynamicPluginServiceServer
// for forward compatibility
type DynamicPluginServiceServer interface {
	Update(context.Context, *DynamicPluginUpdateRequest) (*DynamicPluginUpdateResponse, error)
	mustEmbedUnimplementedDynamicPluginServiceServer()
}

// UnimplementedDynamicPluginServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDynamicPluginServiceServer struct {
}

func (UnimplementedDynamicPluginServiceServer) Update(context.Context, *DynamicPluginUpdateRequest) (*DynamicPluginUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedDynamicPluginServiceServer) mustEmbedUnimplementedDynamicPluginServiceServer() {}

// UnsafeDynamicPluginServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DynamicPluginServiceServer will
// result in compilation errors.
type UnsafeDynamicPluginServiceServer interface {
	mustEmbedUnimplementedDynamicPluginServiceServer()
}

func RegisterDynamicPluginServiceServer(s grpc.ServiceRegistrar, srv DynamicPluginServiceServer) {
	s.RegisterService(&DynamicPluginService_ServiceDesc, srv)
}

func _DynamicPluginService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DynamicPluginUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DynamicPluginServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ioc_golang.aop.dynamic_plugin.DynamicPluginService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DynamicPluginServiceServer).Update(ctx, req.(*DynamicPluginUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DynamicPluginService_ServiceDesc is the grpc.ServiceDesc for DynamicPluginService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DynamicPluginService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ioc_golang.aop.dynamic_plugin.DynamicPluginService",
	HandlerType: (*DynamicPluginServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Update",
			Handler:    _DynamicPluginService_Update_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "extension/aop/dynamic_plugin/api/ioc_golang/aop/dynamic_plugin/dynamic_plugin.proto",
}
