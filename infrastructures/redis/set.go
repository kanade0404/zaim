package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func SetRequestSecret(ctx context.Context, c *redis.Client, token string) error {
	return c.Set(ctx, requestSecret, token, 0).Err()
}
func SetOAuthToken(ctx context.Context, c *redis.Client, token string) error {
	return c.Set(ctx, oauthToken, token, 0).Err()
}
func SetOAuthTokenSecret(ctx context.Context, c *redis.Client, token string) error {
	return c.Set(ctx, oauthTokenSecret, token, 0).Err()
}
