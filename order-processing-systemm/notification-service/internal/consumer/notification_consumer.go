package consumer

import (
	"encoding/json"
	"log"

	"order-processing-system/notification-service/internal/service"
)

type OrderCreatedEvent struct {
	ID int `json:"id"`
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

	var event OrderCreatedEvent

	err := json.Unmarshal(data, &event)

	if err != nil {
		log.Println(err)
		return
	}

	c.Service.SendNotification(event.ID)

	log.Println("Notification sent")
}
