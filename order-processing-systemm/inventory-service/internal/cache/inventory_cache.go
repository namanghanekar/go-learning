package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type InventoryCache struct {
	Client *redis.Client
}

func NewInventoryCache() *InventoryCache {

	client := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
		},
	)

	return &InventoryCache{
		Client: client,
	}
}

func (c *InventoryCache) Ping() error {

	return c.Client.Ping(
		context.Background(),
	).Err()
}
