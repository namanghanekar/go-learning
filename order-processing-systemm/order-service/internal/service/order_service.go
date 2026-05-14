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
func (s *OrderService) UpdateOrderStatus(
	orderID int,
	status string,
) error {

	return s.Repo.UpdateOrderStatus(
		orderID,
		status,
	)
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

	order := &model.Order{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Amount:    req.Amount,

		// IMPORTANT
		Status: "PENDING_PAYMENT",
	}

	fmt.Println("Saving order")

	err = s.Repo.CreateOrder(order)

	if err != nil {
		return nil, err
	}

	fmt.Println("Order saved")

	return order, nil

	err = s.Producer.Publish(order)

	if err != nil {

		order.Status = "FAILED"

		s.Repo.UpdateOrder(order)

		return nil, err
	}

	fmt.Println("Kafka event published")

	return order, nil
}
