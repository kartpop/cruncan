// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package model

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

// PaymentsServiceClient is the client API for PaymentsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PaymentsServiceClient interface {
	DoTransaction(ctx context.Context, in *TransactionRequest, opts ...grpc.CallOption) (*TransactionResponse, error)
}

type paymentsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentsServiceClient(cc grpc.ClientConnInterface) PaymentsServiceClient {
	return &paymentsServiceClient{cc}
}

func (c *paymentsServiceClient) DoTransaction(ctx context.Context, in *TransactionRequest, opts ...grpc.CallOption) (*TransactionResponse, error) {
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, "/model.PaymentsService/DoTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentsServiceServer is the server API for PaymentsService service.
// All implementations must embed UnimplementedPaymentsServiceServer
// for forward compatibility
type PaymentsServiceServer interface {
	DoTransaction(context.Context, *TransactionRequest) (*TransactionResponse, error)
	mustEmbedUnimplementedPaymentsServiceServer()
}

// UnimplementedPaymentsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPaymentsServiceServer struct {
}

func (UnimplementedPaymentsServiceServer) DoTransaction(context.Context, *TransactionRequest) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DoTransaction not implemented")
}
func (UnimplementedPaymentsServiceServer) mustEmbedUnimplementedPaymentsServiceServer() {}

// UnsafePaymentsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PaymentsServiceServer will
// result in compilation errors.
type UnsafePaymentsServiceServer interface {
	mustEmbedUnimplementedPaymentsServiceServer()
}

func RegisterPaymentsServiceServer(s grpc.ServiceRegistrar, srv PaymentsServiceServer) {
	s.RegisterService(&PaymentsService_ServiceDesc, srv)
}

func _PaymentsService_DoTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServiceServer).DoTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/model.PaymentsService/DoTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServiceServer).DoTransaction(ctx, req.(*TransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PaymentsService_ServiceDesc is the grpc.ServiceDesc for PaymentsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PaymentsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "model.PaymentsService",
	HandlerType: (*PaymentsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DoTransaction",
			Handler:    _PaymentsService_DoTransaction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "model/model.proto",
}
