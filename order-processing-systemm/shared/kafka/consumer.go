package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	Reader *kafka.Reader
}

func NewConsumer(
	broker string,
	topic string,
	groupID string,
) *Consumer {

	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: groupID,
		},
	)

	return &Consumer{
		Reader: reader,
	}
}

func (c *Consumer) Consume(
	handler func([]byte),
) {

	for {

		message, err := c.Reader.ReadMessage(
			context.Background(),
		)

		if err != nil {
			log.Println(err)
			continue
		}

		handler(message.Value)
	}
}
