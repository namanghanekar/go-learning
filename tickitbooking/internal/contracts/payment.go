package contracts

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type RegisterHoldRequest struct {
	SeatID      string    `json:"seatId"`
	SeatNumber  string    `json:"seatNumber"`
	RoomID      string    `json:"roomId"`
	HoldToken   string    `json:"holdToken"`
	UserID      string    `json:"userId"`
	ExpiresAt   time.Time `json:"expiresAt"`
	AmountCents int64     `json:"amountCents"`
}

type SubmitPaymentRequest struct {
	SeatID      string `json:"seatId"`
	HoldToken   string `json:"holdToken"`
	UserID      string `json:"userId"`
	AmountCents int64  `json:"amountCents"`
}

type PaymentResponse struct {
	Accepted   bool      `json:"accepted"`
	Message    string    `json:"message"`
	PaymentRef string    `json:"paymentRef"`
	At         time.Time `json:"at"`
}

type PaymentServiceClient interface {
	RegisterHold(ctx context.Context, in *RegisterHoldRequest, opts ...grpc.CallOption) (*PaymentResponse, error)
	SubmitPayment(ctx context.Context, in *SubmitPaymentRequest, opts ...grpc.CallOption) (*PaymentResponse, error)
}

type PaymentServiceServer interface {
	RegisterHold(context.Context, *RegisterHoldRequest) (*PaymentResponse, error)
	SubmitPayment(context.Context, *SubmitPaymentRequest) (*PaymentResponse, error)
}

type paymentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentServiceClient(cc grpc.ClientConnInterface) PaymentServiceClient {
	return &paymentServiceClient{cc: cc}
}

func (c *paymentServiceClient) RegisterHold(ctx context.Context, in *RegisterHoldRequest, opts ...grpc.CallOption) (*PaymentResponse, error) {
	out := new(PaymentResponse)
	err := c.cc.Invoke(ctx, "/payment.PaymentService/RegisterHold", in, out, opts...)
	return out, err
}

func (c *paymentServiceClient) SubmitPayment(ctx context.Context, in *SubmitPaymentRequest, opts ...grpc.CallOption) (*PaymentResponse, error) {
	out := new(PaymentResponse)
	err := c.cc.Invoke(ctx, "/payment.PaymentService/SubmitPayment", in, out, opts...)
	return out, err
}

func RegisterPaymentServiceServer(s grpc.ServiceRegistrar, srv PaymentServiceServer) {
	s.RegisterService(&PaymentService_ServiceDesc, srv)
}

var PaymentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "payment.PaymentService",
	HandlerType: (*PaymentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterHold",
			Handler:    _PaymentService_RegisterHold_Handler,
		},
		{
			MethodName: "SubmitPayment",
			Handler:    _PaymentService_SubmitPayment_Handler,
		},
	},
}

func _PaymentService_RegisterHold_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(RegisterHoldRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).RegisterHold(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/payment.PaymentService/RegisterHold"}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(PaymentServiceServer).RegisterHold(ctx, req.(*RegisterHoldRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_SubmitPayment_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(SubmitPaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).SubmitPayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/payment.PaymentService/SubmitPayment"}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(PaymentServiceServer).SubmitPayment(ctx, req.(*SubmitPaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}
