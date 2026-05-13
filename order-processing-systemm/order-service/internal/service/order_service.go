package service

import (
	"errors"
	"fmt"

	"order-processing-system/order-service/internal/client"
	"order-processing-system/order-service/internal/dto"
	"order-processing-system/order-service/internal/kafka"
	"order-processing-system/order-service/internal/model"
	"order-processing-system/order-service/internal/repository"
)

type OrderService struct {
	Repo            *repository.OrderRepository
	InventoryClient *client.InventoryClient
	PaymentClient   *client.PaymentClient
	Producer        *kafka.Producer
}

func NewOrderService(
	repo *repository.OrderRepository,
	inventoryClient *client.InventoryClient,
	paymentClient *client.PaymentClient,
	producer *kafka.Producer,
) *OrderService {

	return &OrderService{
		Repo:            repo,
		InventoryClient: inventoryClient,
		PaymentClient:   paymentClient,
		Producer:        producer,
	}
}

func (s *OrderService) CreateOrder(
	req dto.CreateOrderRequest,
) (*model.Order, error) {

	fmt.Println("Checking inventory")

	available, err := s.InventoryClient.CheckStock(
		req.ProductID,
		req.Quantity,
	)

	if err != nil {
		return nil, err
	}

	if !available {
		return nil, errors.New("stock not available")
	}

	fmt.Println("Inventory available")

	fmt.Println("Processing payment")

	success, err := s.PaymentClient.ProcessPayment(
		req.Amount,
	)

	if err != nil {
		return nil, err
	}

	if !success {
		return nil, errors.New("payment failed")
	}

	fmt.Println("Payment successful")

	order := &model.Order{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Amount:    req.Amount,
		Status:    "CONFIRMED",
	}

	fmt.Println("Saving order")

	err = s.Repo.CreateOrder(order)

	if err != nil {
		return nil, err
	}

	fmt.Println("Order saved")

	err = s.Producer.Publish(order)

	if err != nil {
		return nil, err
	}

	fmt.Println("Kafka event published")

	return order, nil
}
