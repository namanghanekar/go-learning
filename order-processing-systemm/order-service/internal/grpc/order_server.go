package grpc

import (
	"context"

	"order-processing-system/order-service/internal/dto"
	"order-processing-system/order-service/internal/service"

	pb "order-processing-system/proto/generated"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer

	Service *service.OrderService
}

func NewOrderServer(
	service *service.OrderService,
) *OrderServer {

	return &OrderServer{
		Service: service,
	}
}

func (s *OrderServer) CreateOrder(
	ctx context.Context,
	req *pb.OrderRequest,
) (*pb.OrderResponse, error) {

	orderReq := dto.CreateOrderRequest{
		UserID:    int(req.UserId),
		ProductID: int(req.ProductId),
		Quantity:  int(req.Quantity),
		Amount:    float64(req.Amount),
	}

	_, err := s.Service.CreateOrder(
		orderReq,
	)

	if err != nil {
		return nil, err
	}

	return &pb.OrderResponse{
		Message: "Order created successfully",
	}, nil
}
