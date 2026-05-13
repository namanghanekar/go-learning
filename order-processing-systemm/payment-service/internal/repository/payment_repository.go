package repository

import (
	"order-processing-system/payment-service/internal/model"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	DB *gorm.DB
}

func NewPaymentRepository(
	db *gorm.DB,
) *PaymentRepository {

	return &PaymentRepository{
		DB: db,
	}
}

func (r *PaymentRepository) CreatePayment(
	payment *model.Payment,
) error {

	return r.DB.Create(payment).Error
}
