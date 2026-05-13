package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Reader *kafka.Reader
}

func NewKafkaConsumer(
	broker string,
	topic string,
	groupID string,
) *KafkaConsumer {

	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: groupID,
		},
	)

	return &KafkaConsumer{
		Reader: reader,
	}
}

func (c *KafkaConsumer) Consume(
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
