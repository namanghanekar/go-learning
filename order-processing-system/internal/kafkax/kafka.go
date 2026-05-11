package kafkax

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"

	"order-processing-system/internal/retry"
)

const (
	TopicOrderEvents          = "order.events"
	TopicNotificationRequests = "notification.requests"
)

type Event struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	OrderID   string          `json:"order_id"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}
type Producer struct {
	writer *kafka.Writer
	log    *slog.Logger
}

func NewProducer(brokers []string, log *slog.Logger) *Producer {
	return &Producer{writer: &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}, log: log}
}
func (p *Producer) Publish(ctx context.Context, topic string, key string, event Event) error {
	event.Timestamp = time.Now().UTC()
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return retry.Do(ctx, 3, 100*time.Millisecond, func(ctx context.Context) error {
		return p.writer.WriteMessages(ctx, kafka.Message{Topic: topic, Key: []byte(key), Value: data})
	})
}
func (p *Producer) Close() error { return p.writer.Close() }
func NewReader(brokers []string, groupID string, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})
}
