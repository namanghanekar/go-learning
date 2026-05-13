package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {

	client := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
		},
	)

	client.Ping(context.Background())

	return client
}
