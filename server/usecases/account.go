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

type UpdateAccountUseCase interface {
	UpdateAccounts(ctx context.Context, userName string) error
}

type updateAccountUseCase struct {
	logger echo.Logger
	uow    commands.AccountUnitOfWork
}

func NewUpdateAccountUseCase(db *bun.DB, logger echo.Logger) UpdateAccountUseCase {
	return updateAccountUseCase{
		logger: logger,
		uow:    commands.NewAccountUnitOfWork(db, logger),
	}
}

var _ UpdateAccountUseCase = (*updateAccountUseCase)(nil)

func (a updateAccountUseCase) UpdateAccounts(ctx context.Context, userName string) error {
	a.logger.Infof("UpdateAccounts start.")
	defer a.logger.Infof("UpdateAccounts end.")
	if err := a.uow.Do(ctx, func(ctx context.Context, uowRepositoryManager commands.AccountRepositoryManager) error {
		user, err := uowRepositoryManager.UserRepositoryManager().FindByName(ctx, userName)
		if err != nil {
			return fmt.Errorf("failed to find user: %w", err)
		}
		cfg, err := uowRepositoryManager.ConfigRepositoryManager().FindByUserID(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("failed to find config: %w", err)
		}
		zaimDriver, err := driver.NewZaimDriver(cfg.OAuthConfig.ConsumerKey, cfg.OAuthConfig.ConsumerSecret, cfg.OAuthToken.Token, cfg.OAuthToken.Secret)
		if err != nil {
			return fmt.Errorf("failed to new zaim driver: %w", err)
		}
		res, err := zaimDriver.Get("home/account", nil)
		if err != nil {
			return fmt.Errorf("failed to get account: %w", err)
		}
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %w", err)
		}
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				a.logger.Errorf("failed to close body: %v\n", err)
			}
		}(res.Body)
		var r struct {
			Accounts  []zaim.Account `json:"accounts"`
			Requested int            `json:"requested"`
		}
		if err := json.Unmarshal(b, &r); err != nil {
			return fmt.Errorf("failed to unmarshal json: %w", err)
		}
		a.logger.Infof("accounts: %+v", r.Accounts)
		for _, account := range r.Accounts {
			accountService := uowRepositoryManager.AccountRepositoryManager()
			if err := accountService.Save(ctx, domains.SaveAccountInput{
				ID:     int64(account.ID),
				Name:   account.Name,
				Sort:   account.Sort,
				Active: account.Active,
				UserID: user.ID,
			}); err != nil {
				return fmt.Errorf("failed to save account: %w", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to do: %w", err)
	}
	return nil
}
