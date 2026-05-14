package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"order-processing-system/shipment-service/internal/client"
	"order-processing-system/shipment-service/internal/service"

	"github.com/segmentio/kafka-go"
)

type PaymentEvent struct {
	OrderID int `json:"OrderID"`
}

func StartConsumer(
	shipmentService *service.ShipmentService,
) {

	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "payment.success",
			GroupID: "shipment-group",
		},
	)

	fmt.Println("Shipment consumer started")

	for {

		msg, err := reader.ReadMessage(
			context.Background(),
		)

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("SHIPMENT EVENT RECEIVED")
		fmt.Println(string(msg.Value))

		var event PaymentEvent

		err = json.Unmarshal(
			msg.Value,
			&event,
		)

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(
			"Creating shipment for order:",
			event.OrderID,
		)

		err = shipmentService.CreateShipment(
			event.OrderID,
		)

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Shipment created")

		orderClient := client.NewOrderClient()

		orderClient.UpdateStatus(
			event.OrderID,
			"SHIPPED",
		)

		fmt.Println("Order status updated to SHIPPED")
	}
}
