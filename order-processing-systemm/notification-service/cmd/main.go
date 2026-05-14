package main

import (
	"context"
	"fmt"
	"log"

	"order-processing-system/notification-service/internal/consumer"
	"order-processing-system/notification-service/internal/model"
	"order-processing-system/notification-service/internal/repository"
	"order-processing-system/notification-service/internal/service"

	"order-processing-system/shared/postgres"

	"github.com/segmentio/kafka-go"
)

func main() {

	db := postgres.ConnectDB()

	fmt.Println(
		"PostgreSQL connected",
	)

	err := db.AutoMigrate(
		&model.Notification{},
	)

	if err != nil {
		log.Fatal(err)
	}

	notificationRepo := repository.NewNotificationRepository(
		db,
	)

	notificationService := service.NewNotificationService(
		notificationRepo,
	)

	notificationConsumer := consumer.NewNotificationConsumer(
		notificationService,
	)

	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "payment.success",
			GroupID: "notification-group",
		},
	)

	defer reader.Close()

	fmt.Println(
		"Notification Service Running",
	)

	for {

		msg, err := reader.ReadMessage(
			context.Background(),
		)

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(
			"Notification event received",
		)

		notificationConsumer.HandleMessage(
			msg.Value,
		)
	}
}
