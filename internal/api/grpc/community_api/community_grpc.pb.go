// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: proto/community.proto

package community_api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	CommunityService_CheckAccess_FullMethodName = "/community_api.CommunityService/CheckAccess"
	CommunityService_GetHeader_FullMethodName   = "/community_api.CommunityService/GetHeader"
)

// CommunityServiceClient is the client API for CommunityService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommunityServiceClient interface {
	CheckAccess(ctx context.Context, in *CheckAccessRequest, opts ...grpc.CallOption) (*CheckAccessResponse, error)
	GetHeader(ctx context.Context, in *GetHeaderRequest, opts ...grpc.CallOption) (*GetHeaderResponse, error)
}

type communityServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCommunityServiceClient(cc grpc.ClientConnInterface) CommunityServiceClient {
	return &communityServiceClient{cc}
}

func (c *communityServiceClient) CheckAccess(ctx context.Context, in *CheckAccessRequest, opts ...grpc.CallOption) (*CheckAccessResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckAccessResponse)
	err := c.cc.Invoke(ctx, CommunityService_CheckAccess_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communityServiceClient) GetHeader(ctx context.Context, in *GetHeaderRequest, opts ...grpc.CallOption) (*GetHeaderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetHeaderResponse)
	err := c.cc.Invoke(ctx, CommunityService_GetHeader_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommunityServiceServer is the server API for CommunityService service.
// All implementations must embed UnimplementedCommunityServiceServer
// for forward compatibility.
type CommunityServiceServer interface {
	CheckAccess(context.Context, *CheckAccessRequest) (*CheckAccessResponse, error)
	GetHeader(context.Context, *GetHeaderRequest) (*GetHeaderResponse, error)
	mustEmbedUnimplementedCommunityServiceServer()
}

// UnimplementedCommunityServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCommunityServiceServer struct{}

func (UnimplementedCommunityServiceServer) CheckAccess(context.Context, *CheckAccessRequest) (*CheckAccessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckAccess not implemented")
}
func (UnimplementedCommunityServiceServer) GetHeader(context.Context, *GetHeaderRequest) (*GetHeaderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHeader not implemented")
}
func (UnimplementedCommunityServiceServer) mustEmbedUnimplementedCommunityServiceServer() {}
func (UnimplementedCommunityServiceServer) testEmbeddedByValue()                          {}

// UnsafeCommunityServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommunityServiceServer will
// result in compilation errors.
type UnsafeCommunityServiceServer interface {
	mustEmbedUnimplementedCommunityServiceServer()
}

func RegisterCommunityServiceServer(s grpc.ServiceRegistrar, srv CommunityServiceServer) {
	// If the following call pancis, it indicates UnimplementedCommunityServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CommunityService_ServiceDesc, srv)
}

func _CommunityService_CheckAccess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckAccessRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunityServiceServer).CheckAccess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CommunityService_CheckAccess_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunityServiceServer).CheckAccess(ctx, req.(*CheckAccessRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommunityService_GetHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetHeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunityServiceServer).GetHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CommunityService_GetHeader_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunityServiceServer).GetHeader(ctx, req.(*GetHeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CommunityService_ServiceDesc is the grpc.ServiceDesc for CommunityService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CommunityService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "community_api.CommunityService",
	HandlerType: (*CommunityServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckAccess",
			Handler:    _CommunityService_CheckAccess_Handler,
		},
		{
			MethodName: "GetHeader",
			Handler:    _CommunityService_GetHeader_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/community.proto",
}
