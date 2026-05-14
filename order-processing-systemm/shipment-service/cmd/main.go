package main

import (
	"fmt"

	"order-processing-system/shipment-service/internal/consumer"
	"order-processing-system/shipment-service/internal/model"
	"order-processing-system/shipment-service/internal/repository"
	"order-processing-system/shipment-service/internal/service"

	"order-processing-system/shared/postgres"
)

func main() {

	db := postgres.ConnectDB()

	db.AutoMigrate(&model.Shipment{})

	shipmentRepo := repository.NewShipmentRepository(
		db,
	)

	shipmentService := service.NewShipmentService(
		shipmentRepo,
	)

	fmt.Println(
		"Shipment Service Running",
	)

	consumer.StartConsumer(
		shipmentService,
	)
}
