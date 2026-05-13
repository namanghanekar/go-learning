package service

import (
	"order-processing-system/payment-service/internal/model"
	"order-processing-system/payment-service/internal/repository"
)

type PaymentService struct {
	Repo *repository.PaymentRepository
}

func NewPaymentService(
	repo *repository.PaymentRepository,
) *PaymentService {

	return &PaymentService{
		Repo: repo,
	}
}

func (s *PaymentService) ProcessPayment(
	orderID int,
	amount float64,
) bool {

	payment := &model.Payment{
		OrderID: orderID,
		Amount:  amount,
		Status:  "SUCCESS",
		Method:  "UPI",
	}

	err := s.Repo.CreatePayment(payment)

	if err != nil {
		return false
	}

	return true
}
