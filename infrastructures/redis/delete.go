package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func DeleteRequestSecret(ctx context.Context, c *redis.Client, url string) error {
	return c.Del(ctx, url).Err()
}
