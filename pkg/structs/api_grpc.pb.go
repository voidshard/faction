// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.3
// source: api.proto

package structs

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

// APIClient is the client API for API service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type APIClient interface {
	Worlds(ctx context.Context, in *GetWorldsRequest, opts ...grpc.CallOption) (*GetWorldsResponse, error)
	SetWorld(ctx context.Context, in *SetWorldRequest, opts ...grpc.CallOption) (*SetWorldResponse, error)
	ListWorlds(ctx context.Context, in *ListWorldsRequest, opts ...grpc.CallOption) (*ListWorldsResponse, error)
	DeleteWorld(ctx context.Context, in *DeleteWorldRequest, opts ...grpc.CallOption) (*DeleteWorldResponse, error)
	Actors(ctx context.Context, in *GetActorsRequest, opts ...grpc.CallOption) (*GetActorsResponse, error)
	SetActors(ctx context.Context, in *SetActorsRequest, opts ...grpc.CallOption) (*SetActorsResponse, error)
	ListActors(ctx context.Context, in *ListActorsRequest, opts ...grpc.CallOption) (*ListActorsResponse, error)
	DeleteActor(ctx context.Context, in *DeleteActorRequest, opts ...grpc.CallOption) (*DeleteActorResponse, error)
	// rpc Races(GetRacesRequest) returns (GetRacesResponse);
	// rpc SetRace(SetRaceRequest) returns (SetRaceResponse);
	//
	// rpc Cultures(GetCulturesRequest) returns (GetCulturesResponse);
	// rpc SetCulture(SetCultureRequest) returns (SetCultureResponse);
	Factions(ctx context.Context, in *GetFactionsRequest, opts ...grpc.CallOption) (*GetFactionsResponse, error)
	SetFaction(ctx context.Context, in *SetFactionsRequest, opts ...grpc.CallOption) (*SetFactionsResponse, error)
	ListFactions(ctx context.Context, in *ListFactionsRequest, opts ...grpc.CallOption) (*ListFactionsResponse, error)
	DeleteFaction(ctx context.Context, in *DeleteFactionRequest, opts ...grpc.CallOption) (*DeleteFactionResponse, error)
}

type aPIClient struct {
	cc grpc.ClientConnInterface
}

func NewAPIClient(cc grpc.ClientConnInterface) APIClient {
	return &aPIClient{cc}
}

func (c *aPIClient) Worlds(ctx context.Context, in *GetWorldsRequest, opts ...grpc.CallOption) (*GetWorldsResponse, error) {
	out := new(GetWorldsResponse)
	err := c.cc.Invoke(ctx, "/API/Worlds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) SetWorld(ctx context.Context, in *SetWorldRequest, opts ...grpc.CallOption) (*SetWorldResponse, error) {
	out := new(SetWorldResponse)
	err := c.cc.Invoke(ctx, "/API/SetWorld", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) ListWorlds(ctx context.Context, in *ListWorldsRequest, opts ...grpc.CallOption) (*ListWorldsResponse, error) {
	out := new(ListWorldsResponse)
	err := c.cc.Invoke(ctx, "/API/ListWorlds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) DeleteWorld(ctx context.Context, in *DeleteWorldRequest, opts ...grpc.CallOption) (*DeleteWorldResponse, error) {
	out := new(DeleteWorldResponse)
	err := c.cc.Invoke(ctx, "/API/DeleteWorld", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) Actors(ctx context.Context, in *GetActorsRequest, opts ...grpc.CallOption) (*GetActorsResponse, error) {
	out := new(GetActorsResponse)
	err := c.cc.Invoke(ctx, "/API/Actors", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) SetActors(ctx context.Context, in *SetActorsRequest, opts ...grpc.CallOption) (*SetActorsResponse, error) {
	out := new(SetActorsResponse)
	err := c.cc.Invoke(ctx, "/API/SetActors", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) ListActors(ctx context.Context, in *ListActorsRequest, opts ...grpc.CallOption) (*ListActorsResponse, error) {
	out := new(ListActorsResponse)
	err := c.cc.Invoke(ctx, "/API/ListActors", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) DeleteActor(ctx context.Context, in *DeleteActorRequest, opts ...grpc.CallOption) (*DeleteActorResponse, error) {
	out := new(DeleteActorResponse)
	err := c.cc.Invoke(ctx, "/API/DeleteActor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) Factions(ctx context.Context, in *GetFactionsRequest, opts ...grpc.CallOption) (*GetFactionsResponse, error) {
	out := new(GetFactionsResponse)
	err := c.cc.Invoke(ctx, "/API/Factions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) SetFaction(ctx context.Context, in *SetFactionsRequest, opts ...grpc.CallOption) (*SetFactionsResponse, error) {
	out := new(SetFactionsResponse)
	err := c.cc.Invoke(ctx, "/API/SetFaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) ListFactions(ctx context.Context, in *ListFactionsRequest, opts ...grpc.CallOption) (*ListFactionsResponse, error) {
	out := new(ListFactionsResponse)
	err := c.cc.Invoke(ctx, "/API/ListFactions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) DeleteFaction(ctx context.Context, in *DeleteFactionRequest, opts ...grpc.CallOption) (*DeleteFactionResponse, error) {
	out := new(DeleteFactionResponse)
	err := c.cc.Invoke(ctx, "/API/DeleteFaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// APIServer is the server API for API service.
// All implementations must embed UnimplementedAPIServer
// for forward compatibility
type APIServer interface {
	Worlds(context.Context, *GetWorldsRequest) (*GetWorldsResponse, error)
	SetWorld(context.Context, *SetWorldRequest) (*SetWorldResponse, error)
	ListWorlds(context.Context, *ListWorldsRequest) (*ListWorldsResponse, error)
	DeleteWorld(context.Context, *DeleteWorldRequest) (*DeleteWorldResponse, error)
	Actors(context.Context, *GetActorsRequest) (*GetActorsResponse, error)
	SetActors(context.Context, *SetActorsRequest) (*SetActorsResponse, error)
	ListActors(context.Context, *ListActorsRequest) (*ListActorsResponse, error)
	DeleteActor(context.Context, *DeleteActorRequest) (*DeleteActorResponse, error)
	// rpc Races(GetRacesRequest) returns (GetRacesResponse);
	// rpc SetRace(SetRaceRequest) returns (SetRaceResponse);
	//
	// rpc Cultures(GetCulturesRequest) returns (GetCulturesResponse);
	// rpc SetCulture(SetCultureRequest) returns (SetCultureResponse);
	Factions(context.Context, *GetFactionsRequest) (*GetFactionsResponse, error)
	SetFaction(context.Context, *SetFactionsRequest) (*SetFactionsResponse, error)
	ListFactions(context.Context, *ListFactionsRequest) (*ListFactionsResponse, error)
	DeleteFaction(context.Context, *DeleteFactionRequest) (*DeleteFactionResponse, error)
	mustEmbedUnimplementedAPIServer()
}

// UnimplementedAPIServer must be embedded to have forward compatible implementations.
type UnimplementedAPIServer struct {
}

func (UnimplementedAPIServer) Worlds(context.Context, *GetWorldsRequest) (*GetWorldsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Worlds not implemented")
}
func (UnimplementedAPIServer) SetWorld(context.Context, *SetWorldRequest) (*SetWorldResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetWorld not implemented")
}
func (UnimplementedAPIServer) ListWorlds(context.Context, *ListWorldsRequest) (*ListWorldsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListWorlds not implemented")
}
func (UnimplementedAPIServer) DeleteWorld(context.Context, *DeleteWorldRequest) (*DeleteWorldResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteWorld not implemented")
}
func (UnimplementedAPIServer) Actors(context.Context, *GetActorsRequest) (*GetActorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Actors not implemented")
}
func (UnimplementedAPIServer) SetActors(context.Context, *SetActorsRequest) (*SetActorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetActors not implemented")
}
func (UnimplementedAPIServer) ListActors(context.Context, *ListActorsRequest) (*ListActorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListActors not implemented")
}
func (UnimplementedAPIServer) DeleteActor(context.Context, *DeleteActorRequest) (*DeleteActorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteActor not implemented")
}
func (UnimplementedAPIServer) Factions(context.Context, *GetFactionsRequest) (*GetFactionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Factions not implemented")
}
func (UnimplementedAPIServer) SetFaction(context.Context, *SetFactionsRequest) (*SetFactionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetFaction not implemented")
}
func (UnimplementedAPIServer) ListFactions(context.Context, *ListFactionsRequest) (*ListFactionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListFactions not implemented")
}
func (UnimplementedAPIServer) DeleteFaction(context.Context, *DeleteFactionRequest) (*DeleteFactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFaction not implemented")
}
func (UnimplementedAPIServer) mustEmbedUnimplementedAPIServer() {}

// UnsafeAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to APIServer will
// result in compilation errors.
type UnsafeAPIServer interface {
	mustEmbedUnimplementedAPIServer()
}

func RegisterAPIServer(s grpc.ServiceRegistrar, srv APIServer) {
	s.RegisterService(&API_ServiceDesc, srv)
}

func _API_Worlds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWorldsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).Worlds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API/Worlds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).Worlds(ctx, req.(*GetWorldsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_SetWorld_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetWorldRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).SetWorld(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API/SetWorld",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).SetWorld(ctx, req.(*SetWorldRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_ListWorlds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListWorldsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).ListWorlds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API/ListWorlds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).ListWorlds(ctx, req.(*ListWorldsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_DeleteWorld_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteWorldRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).DeleteWorld(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API/DeleteWorld",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).DeleteWorld(ctx, req.(*DeleteWorldRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_Actors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetActorsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).Actors(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API/Actors",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).Actors(ctx, req.(*GetActorsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_SetActors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetActorsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).SetActors(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API/SetActors",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).SetActors(ctx, req.(*SetActorsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_ListActors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListActorsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).ListActors(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API/ListActors",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).ListActors(ctx, req.(*ListActorsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_DeleteActor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteActorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).DeleteActor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API/DeleteActor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).DeleteActor(ctx, req.(*DeleteActorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_Factions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFactionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).Factions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API/Factions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).Factions(ctx, req.(*GetFactionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_SetFaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetFactionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).SetFaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API/SetFaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).SetFaction(ctx, req.(*SetFactionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_ListFactions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListFactionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).ListFactions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API/ListFactions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).ListFactions(ctx, req.(*ListFactionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_DeleteFaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).DeleteFaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API/DeleteFaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).DeleteFaction(ctx, req.(*DeleteFactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// API_ServiceDesc is the grpc.ServiceDesc for API service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var API_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "API",
	HandlerType: (*APIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Worlds",
			Handler:    _API_Worlds_Handler,
		},
		{
			MethodName: "SetWorld",
			Handler:    _API_SetWorld_Handler,
		},
		{
			MethodName: "ListWorlds",
			Handler:    _API_ListWorlds_Handler,
		},
		{
			MethodName: "DeleteWorld",
			Handler:    _API_DeleteWorld_Handler,
		},
		{
			MethodName: "Actors",
			Handler:    _API_Actors_Handler,
		},
		{
			MethodName: "SetActors",
			Handler:    _API_SetActors_Handler,
		},
		{
			MethodName: "ListActors",
			Handler:    _API_ListActors_Handler,
		},
		{
			MethodName: "DeleteActor",
			Handler:    _API_DeleteActor_Handler,
		},
		{
			MethodName: "Factions",
			Handler:    _API_Factions_Handler,
		},
		{
			MethodName: "SetFaction",
			Handler:    _API_SetFaction_Handler,
		},
		{
			MethodName: "ListFactions",
			Handler:    _API_ListFactions_Handler,
		},
		{
			MethodName: "DeleteFaction",
			Handler:    _API_DeleteFaction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
