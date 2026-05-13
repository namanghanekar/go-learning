package main

import (
	"fmt"

	"order-processing-system/shipment-service/internal/consumer"
)

func main() {

	fmt.Println(
		"Shipment Service Running",
	)

	consumer.StartConsumer()
}
