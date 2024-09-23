package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kanade0404/zaim/server/commands"
	"github.com/kanade0404/zaim/server/domains"
	"github.com/kanade0404/zaim/server/driver"
	"github.com/kanade0404/zaim/server/entity/zaim"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"io"
)

type UpdateCategoryUseCase interface {
	UpdateCategories(ctx context.Context, userName string) error
}
type updateCategoryUseCase struct {
	logger echo.Logger
	uow    commands.CategoryUnitOfWork
}

func (u updateCategoryUseCase) UpdateCategories(ctx context.Context, userName string) error {
	u.logger.Infof("UpdateCategory start. userName: %s", userName)
	defer u.logger.Infof("UpdateCategory end. userName: %s", userName)

	if err := u.uow.Do(ctx, func(ctx context.Context, uowRepositoryManager commands.CategoryRepositoryManager) error {
		user, err := uowRepositoryManager.UserRepositoryManager().FindByName(ctx, userName)
		if err != nil {
			return fmt.Errorf("failed to find user: %v", err)
		}
		cfg, err := uowRepositoryManager.ConfigRepositoryManager().FindByUserID(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("failed to find config: %w", err)
		}
		zaimDriver, err := driver.NewZaimDriver(cfg.OAuthConfig.ConsumerKey, cfg.OAuthConfig.ConsumerSecret, cfg.OAuthToken.Token, cfg.OAuthToken.Secret)
		if err != nil {
			return fmt.Errorf("failed to new zaim driver: %w", err)
		}
		res, err := zaimDriver.Get("home/category", nil)
		if err != nil {
			return fmt.Errorf("failed to get category: %v", err)
		}
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %v", err)
		}
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				u.logger.Infof("failed to close body: %v", err)
			}
		}(res.Body)
		var r struct {
			Categories []zaim.Category `json:"categories"`
			Requested  int             `json:"requested"`
		}
		if err := json.Unmarshal(b, &r); err != nil {
			return fmt.Errorf("failed to unmarshal: %v", err)
		}
		for i := range r.Categories {
			if err := uowRepositoryManager.CategoryRepositoryManager().Save(ctx, domains.SaveCategoryInput{
				ID:     r.Categories[i].ID,
				Name:   r.Categories[i].Name,
				ModeID: r.Categories[i].Mode,
				Sort:   r.Categories[i].Sort,
				Active: r.Categories[i].Active,
				UserID: user.ID,
			}); err != nil {
				return fmt.Errorf("failed to save category: %v", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to do uow: %v", err)
	}
	return nil
}

func NewUpdateCategoryUseCase(db *bun.DB, logger echo.Logger) UpdateCategoryUseCase {
	return &updateCategoryUseCase{
		logger: logger,
		uow:    commands.NewCategoryUnitOfWork(db, logger),
	}
}

var _ UpdateCategoryUseCase = (*updateCategoryUseCase)(nil)
