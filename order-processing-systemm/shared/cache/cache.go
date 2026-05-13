package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Cache struct {
	Client *redis.Client
}

func NewCache() *Cache {

	client := redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
	)

	return &Cache{
		Client: client,
	}
}

func (c *Cache) Set(
	key string,
	value interface{},
	expiration time.Duration,
) error {

	data, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return c.Client.Set(
		ctx,
		key,
		data,
		expiration,
	).Err()
}

func (c *Cache) Get(
	key string,
	dest interface{},
) error {

	val, err := c.Client.Get(
		ctx,
		key,
	).Result()

	if err != nil {
		return err
	}

	return json.Unmarshal(
		[]byte(val),
		dest,
	)
}

func (c *Cache) Delete(
	key string,
) error {

	return c.Client.Del(
		ctx,
		key,
	).Err()
}

func (c *Cache) Exists(
	key string,
) (bool, error) {

	count, err := c.Client.Exists(
		ctx,
		key,
	).Result()

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
