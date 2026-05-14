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
func (r *OrderRepository) UpdateOrder(
	order *model.Order,
) error {

	return r.DB.Save(order).Error
}
func (r *OrderRepository) UpdateOrderStatus(
	orderID int,
	status string,
) error {

	return r.DB.Model(&model.Order{}).
		Where("id = ?", orderID).
		Update("status", status).Error
}
