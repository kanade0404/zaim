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

type GenreUnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context, uowRepositoryManager GenreRepositoryManager) error) error
}
type GenreRepositoryManager interface {
	GenreRepositoryManager() domains.GenreRepository
	UserRepositoryManager() domains.UserRepository
	CategoryRepositoryManager() domains.CategoryRepository
	ConfigRepositoryManager() domains.ConfigRepository
}
type genreRepositoryManager struct {
	genreRepositoryManager    domains.GenreRepository
	userRepositoryManager     domains.UserRepository
	categoryRepositoryManager domains.CategoryRepository
	configRepositoryManager   domains.ConfigRepository
}

func (g genreRepositoryManager) ConfigRepositoryManager() domains.ConfigRepository {
	return g.configRepositoryManager
}

func (g genreRepositoryManager) GenreRepositoryManager() domains.GenreRepository {
	return g.genreRepositoryManager
}

func (g genreRepositoryManager) UserRepositoryManager() domains.UserRepository {
	return g.userRepositoryManager
}

func (g genreRepositoryManager) CategoryRepositoryManager() domains.CategoryRepository {
	return g.categoryRepositoryManager
}

var _ GenreRepositoryManager = (*genreRepositoryManager)(nil)

func newGenreRepositoryManager(repository domains.GenreRepository, userRepository domains.UserRepository, categoryRepository domains.CategoryRepository, configRepository domains.ConfigRepository) GenreRepositoryManager {
	return &genreRepositoryManager{
		genreRepositoryManager:    repository,
		userRepositoryManager:     userRepository,
		categoryRepositoryManager: categoryRepository,
		configRepositoryManager:   configRepository,
	}
}

type genreUnitOfWork struct {
	// TODO: 抽象化する
	db *bun.DB
	// TODO: 抽象化する
	logger echo.Logger
}

func (g *genreUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context, uowRepositoryManager GenreRepositoryManager) error) error {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	repositoryManager := newGenreRepositoryManager(service.NewGenreService(tx, g.logger), service.NewUserService(tx, g.logger), service.NewCategoryService(tx, g.logger), service.NewConfigService(tx, g.logger))
	if err := fn(ctx, repositoryManager); err != nil {
		return driver.Rollback(tx, fmt.Errorf("failed to do: %w", err))
	}
	if err := tx.Commit(); err != nil {
		return driver.Rollback(tx, fmt.Errorf("failed to commit: %w", err))
	}
	return nil
}

func NewGenreUnitOfWork(db *bun.DB, logger echo.Logger) GenreUnitOfWork {
	return &genreUnitOfWork{db: db, logger: logger}
}

var _ GenreUnitOfWork = (*genreUnitOfWork)(nil)
