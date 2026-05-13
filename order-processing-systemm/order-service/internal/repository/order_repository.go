package repository

import (
	"order-processing-system/order-service/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		DB: db,
	}
}

func (r *OrderRepository) CreateOrder(order *model.Order) error {
	return r.DB.Create(order).Error
}
