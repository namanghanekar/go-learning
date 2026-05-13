package grpc

import (
	"context"

	"order-processing-system/payment-service/internal/service"

	pb "order-processing-system/proto/generated"
)

type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer

	Service *service.PaymentService
}

func NewPaymentServer(
	service *service.PaymentService,
) *PaymentServer {

	return &PaymentServer{
		Service: service,
	}
}

func (s *PaymentServer) ProcessPayment(
	ctx context.Context,
	req *pb.PaymentRequest,
) (*pb.PaymentResponse, error) {

	success := s.Service.ProcessPayment(
		int(req.OrderId),
		req.Amount,
	)

	return &pb.PaymentResponse{
		Success: success,
	}, nil
}
