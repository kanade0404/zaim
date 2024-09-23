package domains

import (
	"context"
	"github.com/dghubble/oauth1"
)

type OAuthToken struct {
	Token  string
	Secret string
}
type Config struct {
	OAuthConfig *oauth1.Config
	OAuthToken
}

type ConfigRepository interface {
	FindByUserID(ctx context.Context, userID int64) (*Config, error)
}
