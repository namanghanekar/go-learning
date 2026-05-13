package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer(
	broker string,
	topic string,
) *Producer {

	writer := &kafka.Writer{
		Addr:  kafka.TCP(broker),
		Topic: topic,
	}

	return &Producer{
		Writer: writer,
	}
}

func (p *Producer) Publish(
	message interface{},
) error {

	data, err := json.Marshal(message)

	if err != nil {
		return err
	}

	return p.Writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: data,
		},
	)
}
