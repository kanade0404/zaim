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

type AccountUnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context, uowRepositoryManager AccountRepositoryManager) error) error
}

type AccountRepositoryManager interface {
	AccountRepositoryManager() domains.AccountRepository
	UserRepositoryManager() domains.UserRepository
	ConfigRepositoryManager() domains.ConfigRepository
}

type accountRepositoryManager struct {
	accountRepositoryManager domains.AccountRepository
	userRepositoryManager    domains.UserRepository
	configRepositoryManager  domains.ConfigRepository
}

func (a accountRepositoryManager) ConfigRepositoryManager() domains.ConfigRepository {
	return a.configRepositoryManager
}

func (a accountRepositoryManager) AccountRepositoryManager() domains.AccountRepository {
	return a.accountRepositoryManager
}

func (a accountRepositoryManager) UserRepositoryManager() domains.UserRepository {
	return a.userRepositoryManager
}

var _ AccountRepositoryManager = (*accountRepositoryManager)(nil)

func newAccountRepositoryManager(accountRepository domains.AccountRepository, userRepository domains.UserRepository, configRepository domains.ConfigRepository) AccountRepositoryManager {
	return &accountRepositoryManager{
		accountRepositoryManager: accountRepository,
		userRepositoryManager:    userRepository,
		configRepositoryManager:  configRepository,
	}
}

type accountUnitOfWork struct {
	// TODO: 抽象化する
	db *bun.DB
	// TODO: 抽象化する
	logger echo.Logger
}

func NewAccountUnitOfWork(db *bun.DB, logger echo.Logger) AccountUnitOfWork {
	return &accountUnitOfWork{db: db, logger: logger}
}

func (u *accountUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context, uowRepositoryManager AccountRepositoryManager) error) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	repositoryManager := newAccountRepositoryManager(service.NewAccountService(tx, u.logger), service.NewUserService(tx, u.logger), service.NewConfigService(tx, u.logger))
	if err := fn(ctx, repositoryManager); err != nil {
		return driver.Rollback(tx, fmt.Errorf("failed to do: %w", err))
	}
	if err := tx.Commit(); err != nil {
		return driver.Rollback(tx, fmt.Errorf("failed to commit: %w", err))
	}
	return nil
}

var _ AccountUnitOfWork = (*accountUnitOfWork)(nil)
