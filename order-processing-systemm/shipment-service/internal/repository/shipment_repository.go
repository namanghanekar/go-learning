package repository

import (
	"order-processing-system/shipment-service/internal/model"

	"gorm.io/gorm"
)

type ShipmentRepository struct {
	DB *gorm.DB
}

func NewShipmentRepository(db *gorm.DB) *ShipmentRepository {

	return &ShipmentRepository{
		DB: db,
	}
}

func (r *ShipmentRepository) CreateShipment(
	shipment *model.Shipment,
) error {

	return r.DB.Create(shipment).Error
}
