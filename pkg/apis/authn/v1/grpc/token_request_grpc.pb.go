// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.1
// source: pkg/apis/authn/v1/grpc/token_request.proto

package grpc

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
	Authentication_CreateToken_FullMethodName = "/yhlooo.com.scaf.authn.v1.Authentication/CreateToken"
)

// AuthenticationClient is the client API for Authentication service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 认证服务
type AuthenticationClient interface {
	CreateToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*TokenRequest, error)
}

type authenticationClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthenticationClient(cc grpc.ClientConnInterface) AuthenticationClient {
	return &authenticationClient{cc}
}

func (c *authenticationClient) CreateToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*TokenRequest, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TokenRequest)
	err := c.cc.Invoke(ctx, Authentication_CreateToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthenticationServer is the server API for Authentication service.
// All implementations must embed UnimplementedAuthenticationServer
// for forward compatibility.
//
// 认证服务
type AuthenticationServer interface {
	CreateToken(context.Context, *TokenRequest) (*TokenRequest, error)
	mustEmbedUnimplementedAuthenticationServer()
}

// UnimplementedAuthenticationServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAuthenticationServer struct{}

func (UnimplementedAuthenticationServer) CreateToken(context.Context, *TokenRequest) (*TokenRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateToken not implemented")
}
func (UnimplementedAuthenticationServer) mustEmbedUnimplementedAuthenticationServer() {}
func (UnimplementedAuthenticationServer) testEmbeddedByValue()                        {}

// UnsafeAuthenticationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthenticationServer will
// result in compilation errors.
type UnsafeAuthenticationServer interface {
	mustEmbedUnimplementedAuthenticationServer()
}

func RegisterAuthenticationServer(s grpc.ServiceRegistrar, srv AuthenticationServer) {
	// If the following call pancis, it indicates UnimplementedAuthenticationServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Authentication_ServiceDesc, srv)
}

func _Authentication_CreateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServer).CreateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Authentication_CreateToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServer).CreateToken(ctx, req.(*TokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Authentication_ServiceDesc is the grpc.ServiceDesc for Authentication service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Authentication_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "yhlooo.com.scaf.authn.v1.Authentication",
	HandlerType: (*AuthenticationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateToken",
			Handler:    _Authentication_CreateToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/apis/authn/v1/grpc/token_request.proto",
}
