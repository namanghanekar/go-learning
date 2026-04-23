package utils

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		panic("Redis not connected")
	}
}
