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
	req *pb.CreateOrderRequest,
) (*pb.CreateOrderResponse, error) {

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

	return &pb.CreateOrderResponse{
		Message: "Order created successfully",
	}, nil
}
func (s *OrderServer) UpdateOrderStatus(
	ctx context.Context,
	req *pb.UpdateOrderStatusRequest,
) (*pb.UpdateOrderStatusResponse, error) {

	err := s.Service.UpdateOrderStatus(
		int(req.OrderId),
		req.Status,
	)

	if err != nil {
		return nil, err
	}

	return &pb.UpdateOrderStatusResponse{
		Success: true,
	}, nil
}
