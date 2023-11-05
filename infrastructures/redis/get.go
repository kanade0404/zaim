package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func GetRequestSecret(ctx context.Context, c *redis.Client) (string, error) {
	return c.Get(ctx, requestSecret).Result()
}
