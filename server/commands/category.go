package commands

import (
	"context"
	"fmt"
	"github.com/kanade0404/zaim/server/domains"
	"github.com/kanade0404/zaim/server/driver"
	"github.com/kanade0404/zaim/server/infrastructures/service"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type CategoryUnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context, uowRepositoryManager CategoryRepositoryManager) error) error
}

type CategoryRepositoryManager interface {
	CategoryRepositoryManager() domains.CategoryRepository
	UserRepositoryManager() domains.UserRepository
	ConfigRepositoryManager() domains.ConfigRepository
}

type categoryRepositoryManager struct {
	categoryRepositoryManager domains.CategoryRepository
	userRepositoryManager     domains.UserRepository
	configRepositoryManager   domains.ConfigRepository
}

func (c categoryRepositoryManager) ConfigRepositoryManager() domains.ConfigRepository {
	return c.configRepositoryManager
}

func (c categoryRepositoryManager) CategoryRepositoryManager() domains.CategoryRepository {
	return c.categoryRepositoryManager
}

func (c categoryRepositoryManager) UserRepositoryManager() domains.UserRepository {
	return c.userRepositoryManager
}

var _ CategoryRepositoryManager = (*categoryRepositoryManager)(nil)

func newCategoryRepositoryManager(repository domains.CategoryRepository, userRepository domains.UserRepository, configRepository domains.ConfigRepository) CategoryRepositoryManager {
	return &categoryRepositoryManager{
		categoryRepositoryManager: repository,
		userRepositoryManager:     userRepository,
		configRepositoryManager:   configRepository,
	}
}

type categoryUnitOfWork struct {
	// TODO: 抽象化する
	db *bun.DB
	// TODO: 抽象化する
	logger echo.Logger
}

func (c *categoryUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context, uowRepositoryManager CategoryRepositoryManager) error) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	repositoryManager := newCategoryRepositoryManager(service.NewCategoryService(tx, c.logger), service.NewUserService(tx, c.logger), service.NewConfigService(tx, c.logger))
	if err := fn(ctx, repositoryManager); err != nil {
		return driver.Rollback(tx, fmt.Errorf("failed to do: %w", err))
	}
	if err := tx.Commit(); err != nil {
		return driver.Rollback(tx, fmt.Errorf("failed to commit: %w", err))
	}
	return nil
}

func NewCategoryUnitOfWork(db *bun.DB, logger echo.Logger) CategoryUnitOfWork {
	return &categoryUnitOfWork{db: db, logger: logger}
}

var _ CategoryUnitOfWork = (*categoryUnitOfWork)(nil)
