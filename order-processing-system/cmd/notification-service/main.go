package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"order-processing-system/internal/config"
	"order-processing-system/internal/kafkax"
	"order-processing-system/internal/logx"
	"order-processing-system/internal/redisx"
	"order-processing-system/internal/retry"
	"order-processing-system/internal/shutdown"
	"order-processing-system/internal/workerpool"
)

type notifier struct {
	redis *redis.Client
	log   *slog.Logger
}

func (n notifier) handle(ctx context.Context, event kafkax.Event) error {
	key := "notification:processed:" + event.ID
	claimed, err := n.redis.SetNX(ctx, key, "processing", 24*time.Hour).Result()
	if err != nil {
		return err
	}
	if !claimed {
		n.log.Info("duplicate notification skipped", "event_id", event.ID, "order_id", event.OrderID)
		return nil
	}
	return retry.Do(ctx, 3, 200*time.Millisecond, func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
			n.log.Info("notification sent", "event_id", event.ID, "order_id", event.OrderID, "type", event.Type)
			return n.redis.Set(ctx, key, "completed", 24*time.Hour).Err()
		}
	})
}

func main() {
	cfg := config.Load("notification-service", ":0")
	log := logx.New(cfg.ServiceName)
	ctx := shutdown.Context()
	rdb, err := redisx.Connect(ctx, cfg.RedisAddr)
	if err != nil {
		log.Error("redis connection failed", "error", err)
		return
	}
	reader := kafkax.NewReader(cfg.KafkaBrokers, "notification-service", kafkax.TopicNotificationRequests)
	defer reader.Close()
	workers := config.Int("NOTIFICATION_WORKERS", 4)
	pool := workerpool.New(workers, 100, log)
	pool.Start(ctx, workers)
	defer pool.Shutdown()
	n := notifier{redis: rdb, log: log}
	log.Info("notification consumer started", "workers", workers)
	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Error("kafka fetch failed", "error", err)
			continue
		}
		var event kafkax.Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Error("invalid event", "error", err)
			_ = reader.CommitMessages(ctx, msg)
			continue
		}
		eventCopy := event
		err = pool.Submit(ctx, workerpool.Job{ID: eventCopy.ID, Payload: msg.Value, Handle: func(ctx context.Context) error {
			return n.handle(ctx, eventCopy)
		}})
		if err != nil {
			log.Error("job submit failed", "error", err)
			continue
		}
		_ = reader.CommitMessages(ctx, msg)
	}
}
