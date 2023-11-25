package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

func SetZaimSecret(ctx context.Context, c *redis.Client, oauthToken, user, secret string) error {
	b, err := json.Marshal(RequestSecret{
		Secret: secret,
		User:   user,
	})
	if err != nil {
		return err
	}
	return c.Set(ctx, oauthToken, b, time.Minute*5).Err()
}
func SetOAuthTokens(ctx context.Context, c *redis.Client, user, token, tokenSecret string) error {
	b, err := json.Marshal(OauthToken{
		Token:  token,
		Secret: tokenSecret,
	})
	if err != nil {
		return err
	}
	return c.Set(ctx, user, b, 0).Err()

}
