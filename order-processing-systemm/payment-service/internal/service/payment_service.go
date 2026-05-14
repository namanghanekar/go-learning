package service

import (
	"fmt"

	"order-processing-system/payment-service/internal/client"
	"order-processing-system/payment-service/internal/kafka"
	"order-processing-system/payment-service/internal/model"
	"order-processing-system/payment-service/internal/repository"
)

type PaymentSuccessEvent struct {
	OrderID int `json:"order_id"`
}

type PaymentService struct {
	Repo        *repository.PaymentRepository
	OrderClient *client.OrderClient
	Producer    *kafka.Producer
}

func NewPaymentService(
	repo *repository.PaymentRepository,
	orderClient *client.OrderClient,
	producer *kafka.Producer,
) *PaymentService {

	return &PaymentService{
		Repo:        repo,
		OrderClient: orderClient,
		Producer:    producer,
	}
}

func (s *PaymentService) ProcessPayment(
	orderID int,
	amount float64,
) bool {

	fmt.Println("Processing payment")

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

	fmt.Println("Payment saved")

	// UPDATE ORDER STATUS
	s.OrderClient.UpdateStatus(
		orderID,
		"PAID",
	)

	fmt.Println("Order status updated")

	// CLEAN EVENT
	event := PaymentSuccessEvent{
		OrderID: orderID,
	}

	// PUBLISH EVENT
	err = s.Producer.Publish(event)

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("Kafka event published")

	return true
}
