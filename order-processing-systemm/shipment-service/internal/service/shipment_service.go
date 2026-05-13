package service

import (
	"fmt"
	"order-processing-system/shipment-service/internal/model"
	"order-processing-system/shipment-service/internal/repository"
)

type ShipmentService struct {
	Repo *repository.ShipmentRepository
}

func NewShipmentService(
	repo *repository.ShipmentRepository,
) *ShipmentService {

	return &ShipmentService{
		Repo: repo,
	}
}

func (s *ShipmentService) CreateShipment(
	orderID int,
) error {

	shipment := &model.Shipment{
		OrderID:    orderID,
		TrackingID: fmt.Sprintf("TRK-%d", orderID),
		Status:     "CREATED",
	}

	return s.Repo.CreateShipment(shipment)
}
