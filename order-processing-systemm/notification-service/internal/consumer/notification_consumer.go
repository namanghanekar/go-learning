package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	"order-processing-system/notification-service/internal/client"
	"order-processing-system/notification-service/internal/service"
)

type PaymentEvent struct {
	OrderID int `json:"OrderID"`
}

type NotificationConsumer struct {
	Service *service.NotificationService
}

func NewNotificationConsumer(
	service *service.NotificationService,
) *NotificationConsumer {

	return &NotificationConsumer{
		Service: service,
	}
}

func (c *NotificationConsumer) HandleMessage(
	data []byte,
) {

	fmt.Println("NOTIFICATION EVENT RECEIVED")
	fmt.Println(string(data))

	var event PaymentEvent

	err := json.Unmarshal(
		data,
		&event,
	)

	if err != nil {
		log.Println(err)
		return
	}

	err = c.Service.SendNotification(
		event.OrderID,
		"Order Completed Successfully",
	)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Notification sent")

	orderClient := client.NewOrderClient()

	orderClient.UpdateStatus(
		event.OrderID,
		"COMPLETED",
	)

	fmt.Println("Order status updated to COMPLETED")
}
