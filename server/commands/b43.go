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

type B43UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context, uowRepositoryManager B43RepositoryManager) error) error
}

type B43RepositoryManager interface {
	AccountRepositoryManager() domains.AccountRepository
	UserRepositoryManager() domains.UserRepository
	CategoryRepository() domains.CategoryRepository
	GenreRepository() domains.GenreRepository
	ConfigRepositoryManager() domains.ConfigRepository
}

type b43RepositoryManager struct {
	accountRepository  domains.AccountRepository
	userRepository     domains.UserRepository
	categoryRepository domains.CategoryRepository
	genreRepository    domains.GenreRepository
	configRepository   domains.ConfigRepository
}

func (a b43RepositoryManager) ConfigRepositoryManager() domains.ConfigRepository {
	return a.configRepository
}

func (a b43RepositoryManager) CategoryRepository() domains.CategoryRepository {
	return a.categoryRepository
}

func (a b43RepositoryManager) GenreRepository() domains.GenreRepository {
	return a.genreRepository
}

func (a b43RepositoryManager) AccountRepositoryManager() domains.AccountRepository {
	return a.accountRepository
}

func (a b43RepositoryManager) UserRepositoryManager() domains.UserRepository {
	return a.userRepository
}

var _ B43RepositoryManager = (*b43RepositoryManager)(nil)

func newB43RepositoryManager(accountRepository domains.AccountRepository, userRepository domains.UserRepository, repository domains.GenreRepository, categoryRepository domains.CategoryRepository, configRepository domains.ConfigRepository) B43RepositoryManager {
	return &b43RepositoryManager{
		accountRepository:  accountRepository,
		userRepository:     userRepository,
		genreRepository:    repository,
		categoryRepository: categoryRepository,
		configRepository:   configRepository,
	}
}

type b43UnitOfWork struct {
	// TODO: 抽象化する
	db *bun.DB
	// TODO: 抽象化する
	logger echo.Logger
}

func NewB43UnitOfWork(db *bun.DB, logger echo.Logger) B43UnitOfWork {
	return &b43UnitOfWork{db: db, logger: logger}
}

func (u *b43UnitOfWork) Do(ctx context.Context, fn func(ctx context.Context, uowRepositoryManager B43RepositoryManager) error) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	repositoryManager := newB43RepositoryManager(service.NewAccountService(tx, u.logger), service.NewUserService(tx, u.logger), service.NewGenreService(tx, u.logger), service.NewCategoryService(tx, u.logger), service.NewConfigService(tx, u.logger))
	if err := fn(ctx, repositoryManager); err != nil {
		return driver.Rollback(tx, fmt.Errorf("failed to do: %w", err))
	}
	if err := tx.Commit(); err != nil {
		return driver.Rollback(tx, fmt.Errorf("failed to commit: %w", err))
	}
	return nil
}
