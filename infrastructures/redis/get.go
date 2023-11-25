package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

func GetRequestSecret(ctx context.Context, c *redis.Client, url string) (RequestSecret, error) {
	result, err := c.Get(ctx, url).Result()
	if err != nil {
		return RequestSecret{}, err
	}
	var requestSecret RequestSecret
	if err := json.Unmarshal([]byte(result), &requestSecret); err != nil {
		return RequestSecret{}, err
	}
	return requestSecret, nil
}

func GetOauthToken(ctx context.Context, c *redis.Client, user string) (OauthToken, error) {
	result, err := c.Get(ctx, user).Result()
	if err != nil {
		return OauthToken{}, err
	}
	var oauthToken OauthToken
	if err := json.Unmarshal([]byte(result), &oauthToken); err != nil {
		return OauthToken{}, err
	}
	return oauthToken, nil
}
