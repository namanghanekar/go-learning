package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Order struct {
	ID int `json:"id"`
}

func StartConsumer() {

	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "order.created",
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

		var order Order

		err = json.Unmarshal(
			msg.Value,
			&order,
		)

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(
			"Order event received",
			order.ID,
		)

		fmt.Println(
			"Shipment created",
		)
	}
}
