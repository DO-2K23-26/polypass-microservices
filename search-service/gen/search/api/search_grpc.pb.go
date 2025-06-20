// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: proto/search.proto

package api

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
	SearchService_SearchFolders_FullMethodName     = "/search.SearchService/SearchFolders"
	SearchService_SearchTags_FullMethodName        = "/search.SearchService/SearchTags"
	SearchService_SearchCredentials_FullMethodName = "/search.SearchService/SearchCredentials"
)

// SearchServiceClient is the client API for SearchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SearchServiceClient interface {
	// Search for folders
	SearchFolders(ctx context.Context, in *SearchFoldersRequest, opts ...grpc.CallOption) (*SearchFoldersResponse, error)
	// Search for tags
	SearchTags(ctx context.Context, in *SearchTagsRequest, opts ...grpc.CallOption) (*SearchTagsResponse, error)
	// Search for credentials
	SearchCredentials(ctx context.Context, in *SearchCredentialsRequest, opts ...grpc.CallOption) (*SearchCredentialsResponse, error)
}

type searchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSearchServiceClient(cc grpc.ClientConnInterface) SearchServiceClient {
	return &searchServiceClient{cc}
}

func (c *searchServiceClient) SearchFolders(ctx context.Context, in *SearchFoldersRequest, opts ...grpc.CallOption) (*SearchFoldersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchFoldersResponse)
	err := c.cc.Invoke(ctx, SearchService_SearchFolders_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) SearchTags(ctx context.Context, in *SearchTagsRequest, opts ...grpc.CallOption) (*SearchTagsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchTagsResponse)
	err := c.cc.Invoke(ctx, SearchService_SearchTags_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) SearchCredentials(ctx context.Context, in *SearchCredentialsRequest, opts ...grpc.CallOption) (*SearchCredentialsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchCredentialsResponse)
	err := c.cc.Invoke(ctx, SearchService_SearchCredentials_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SearchServiceServer is the server API for SearchService service.
// All implementations must embed UnimplementedSearchServiceServer
// for forward compatibility.
type SearchServiceServer interface {
	// Search for folders
	SearchFolders(context.Context, *SearchFoldersRequest) (*SearchFoldersResponse, error)
	// Search for tags
	SearchTags(context.Context, *SearchTagsRequest) (*SearchTagsResponse, error)
	// Search for credentials
	SearchCredentials(context.Context, *SearchCredentialsRequest) (*SearchCredentialsResponse, error)
	mustEmbedUnimplementedSearchServiceServer()
}

// UnimplementedSearchServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSearchServiceServer struct{}

func (UnimplementedSearchServiceServer) SearchFolders(context.Context, *SearchFoldersRequest) (*SearchFoldersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchFolders not implemented")
}
func (UnimplementedSearchServiceServer) SearchTags(context.Context, *SearchTagsRequest) (*SearchTagsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchTags not implemented")
}
func (UnimplementedSearchServiceServer) SearchCredentials(context.Context, *SearchCredentialsRequest) (*SearchCredentialsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchCredentials not implemented")
}
func (UnimplementedSearchServiceServer) mustEmbedUnimplementedSearchServiceServer() {}
func (UnimplementedSearchServiceServer) testEmbeddedByValue()                       {}

// UnsafeSearchServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SearchServiceServer will
// result in compilation errors.
type UnsafeSearchServiceServer interface {
	mustEmbedUnimplementedSearchServiceServer()
}

func RegisterSearchServiceServer(s grpc.ServiceRegistrar, srv SearchServiceServer) {
	// If the following call pancis, it indicates UnimplementedSearchServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SearchService_ServiceDesc, srv)
}

func _SearchService_SearchFolders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchFoldersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).SearchFolders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SearchService_SearchFolders_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).SearchFolders(ctx, req.(*SearchFoldersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_SearchTags_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchTagsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).SearchTags(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SearchService_SearchTags_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).SearchTags(ctx, req.(*SearchTagsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_SearchCredentials_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchCredentialsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).SearchCredentials(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SearchService_SearchCredentials_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).SearchCredentials(ctx, req.(*SearchCredentialsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SearchService_ServiceDesc is the grpc.ServiceDesc for SearchService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SearchService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "search.SearchService",
	HandlerType: (*SearchServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SearchFolders",
			Handler:    _SearchService_SearchFolders_Handler,
		},
		{
			MethodName: "SearchTags",
			Handler:    _SearchService_SearchTags_Handler,
		},
		{
			MethodName: "SearchCredentials",
			Handler:    _SearchService_SearchCredentials_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/search.proto",
}
