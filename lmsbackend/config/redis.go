package config

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var Ctx = context.Background()

func ConnectRedis() {

	RDB = redis.NewClient(&redis.Options{
		Addr: AppConfig.RedisAddr,
	})

	// test connection
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		panic("❌ Redis connection failed: " + err.Error())
	}

	fmt.Println("✅ Redis connected")
}
