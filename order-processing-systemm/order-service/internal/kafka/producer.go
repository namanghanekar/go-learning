package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer() *Producer {

	writer := &kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"),
		Topic: "order.created",
	}

	return &Producer{
		Writer: writer,
	}
}

func (p *Producer) Publish(
	data interface{},
) error {

	bytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	err = p.Writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: bytes,
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("Kafka event published")

	return nil
}
