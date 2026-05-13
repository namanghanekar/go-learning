package main

import (
	"fmt"

	"order-processing-system/notification-service/internal/consumer"
	"order-processing-system/notification-service/internal/email"
	"order-processing-system/notification-service/internal/kafka"
	"order-processing-system/notification-service/internal/service"
	"order-processing-system/notification-service/internal/sms"

	"github.com/gin-gonic/gin"
)

func main() {

	emailSender := email.NewEmailSender()

	smsSender := sms.NewSMSSender()

	notificationService := service.NewNotificationService(
		emailSender,
		smsSender,
	)

	notificationConsumer := consumer.NewNotificationConsumer(
		notificationService,
	)

	kafkaConsumer := kafka.NewKafkaConsumer(
		"localhost:9092",
		"order.created",
		"notification-group",
	)

	go kafkaConsumer.Consume(
		notificationConsumer.HandleMessage,
	)

	router := gin.Default()

	router.GET(
		"/health",
		func(c *gin.Context) {

			c.JSON(200, gin.H{
				"message": "Notification Service Running",
			})
		},
	)

	fmt.Println(
		"Notification Service Running On Port 8084",
	)

	router.Run(":8084")
}
