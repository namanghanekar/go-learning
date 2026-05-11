package redisx

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

func Connect(ctx context.Context, addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{Addr: strings.TrimSpace(addr)})
	return client, client.Ping(ctx).Err()
}
