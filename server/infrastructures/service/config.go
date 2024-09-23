package service

import (
	"context"
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/kanade0404/zaim/server/domains"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type ConfigService struct {
	baseService
}

func NewConfigService(tx bun.Tx, logger echo.Logger) ConfigService {
	return ConfigService{
		baseService: baseService{db: tx, logger: logger},
	}

}

func (c ConfigService) FindByUserID(ctx context.Context, userID int64) (*domains.Config, error) {
	var (
		zaimApp   domains.ZaimApplication
		zaimOauth domains.ZaimOAuth
	)
	if err := c.db.NewSelect().Model(&zaimApp).
		Relation("EnableZaimApplicationEvents", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Order("enabled DESC")
		}).
		Where("user_id = ?", userID).
		Limit(1).
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("failed to select zaim application: %v", err)
	}
	if err := c.db.NewSelect().
		NewSelect().
		Model(&zaimOauth).
		Relation("EnableZaimOAuthEvent", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Order("enabled DESC")
		}).
		Where("zaim_app_id = ?", zaimApp.ID).
		Limit(1).
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("failed to select zaim oauth: %v", err)
	}
	config := &domains.Config{
		OAuthConfig: &oauth1.Config{
			ConsumerKey:    zaimApp.ConsumerKey,
			ConsumerSecret: zaimApp.ConsumerSecret,
		},
		OAuthToken: domains.OAuthToken{
			Token:  zaimOauth.Token,
			Secret: zaimOauth.Secret,
		},
	}
	return config, nil
}

var _ domains.ConfigRepository = (*ConfigService)(nil)
