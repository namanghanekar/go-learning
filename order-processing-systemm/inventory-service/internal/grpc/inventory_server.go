package grpc

import (
	"context"

	"order-processing-system/inventory-service/internal/service"

	pb "order-processing-system/proto/generated"
)

type InventoryServer struct {
	pb.UnimplementedInventoryServiceServer

	Service *service.InventoryService
}

func NewInventoryServer(
	service *service.InventoryService,
) *InventoryServer {

	return &InventoryServer{
		Service: service,
	}
}

func (s *InventoryServer) CheckStock(
	ctx context.Context,
	req *pb.InventoryRequest,
) (*pb.InventoryResponse, error) {

	available := s.Service.CheckStock(
		int(req.ProductId),
		int(req.Quantity),
	)

	return &pb.InventoryResponse{
		Available: available,
	}, nil
}
