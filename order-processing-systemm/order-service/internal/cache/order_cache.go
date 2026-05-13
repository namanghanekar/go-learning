package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type OrderCache struct {
	Client *redis.Client
}

func NewOrderCache() *OrderCache {

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &OrderCache{
		Client: client,
	}
}

func (c *OrderCache) Set(
	key string,
	value interface{},
) error {

	ctx := context.Background()

	data, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return c.Client.Set(
		ctx,
		key,
		data,
		0,
	).Err()
}

func (c *OrderCache) Get(
	key string,
) (string, error) {

	ctx := context.Background()

	val, err := c.Client.Get(
		ctx,
		key,
	).Result()

	if err != nil {
		return "", err
	}

	fmt.Println("Cache hit")

	return val, nil
}
