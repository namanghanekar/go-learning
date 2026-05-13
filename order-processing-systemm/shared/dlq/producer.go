package dlq

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewDLQProducer(topic string) *Producer {

	writer := &kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"),
		Topic: topic,
	}

	return &Producer{
		Writer: writer,
	}
}

func (p *Producer) Publish(data interface{}) error {

	bytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	return p.Writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: bytes,
		},
	)
}
