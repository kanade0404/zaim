package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kanade0404/zaim/server/domains"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type AccountService struct {
	baseService
}

func (a AccountService) FindB43(ctx context.Context, input domains.FindB43Input) (*domains.Account, error) {
	a.logger.Infof("Find b43 account start. input: %v", input)
	defer a.logger.Infof("Find b43 account end. input: %v", input)
	var result struct {
		ID        uuid.UUID `bun:"id"`
		AccountID int64     `bun:"account_id"`
		UserID    int64     `bun:"user_id"`
	}
	subQuery := a.db.NewSelect().
		Model((*domains.Account)(nil)).
		ColumnExpr("a.id,a.account_id,user_id,b43_account.modified,max(activated) AS last_activated,max(inactivated) AS last_inactivated").
		Join("INNER JOIN account_modified_event ON account_modified_event.account_id = a.id").
		Join("INNER JOIN b43_account ON b43_account.account_id = a.id").
		Join("LEFT JOIN active_account ON active_account.account_id = a.id").
		Join("LEFT JOIN inactive_account ON inactive_account.account_id = a.id").
		Where("user_id = ?", input.UserID).
		Group("a.id", "a.account_id", "user_id", "b43_account.modified")
	q := a.db.NewSelect().TableExpr("(?) AS accounts", subQuery)
	if !input.IncludeInActive {
		q = q.Where("accounts.last_activated IS NOT NULL").
			WhereGroup(" AND ", func(query *bun.SelectQuery) *bun.SelectQuery {
				return query.Where("accounts.last_inactivated IS NULL").
					WhereOr("accounts.last_activated > accounts.last_inactivated")
			})
	}
	if err := q.Column("id", "account_id", "user_id").
		Order("modified DESC", "last_activated DESC", "last_inactivated DESC").
		Limit(1).
		Scan(ctx, &result); err != nil {
		return nil, fmt.Errorf("failed to select account: %v", err)
	}
	account := &domains.Account{
		ID:        result.ID,
		AccountID: result.AccountID,
		UserID:    result.UserID,
	}
	return account, nil
}

func (a AccountService) FindByName(ctx context.Context, input domains.FindByNameAccountInput) (*domains.Account, error) {
	a.logger.Infof("Find account start. input: %v", input)
	defer a.logger.Infof("Find account end. input: %v", input)
	var result struct {
		ID        uuid.UUID `bun:"id"`
		AccountID int64     `bun:"account_id"`
		UserID    int64     `bun:"user_id"`
	}
	subQuery := a.db.NewSelect().
		Model((*domains.Account)(nil)).
		ColumnExpr("a.id,a.account_id,user_id,max(activated) as last_activated,max(inactivated) as last_inactivated").
		Join("INNER JOIN account_modified_event ON account_modified_event.account_id = a.id").
		Join("LEFT JOIN active_account ON active_account.account_id = a.id").
		Join("LEFT JOIN inactive_account ON inactive_account.account_id = a.id").
		Where("user_id = ?", input.UserID).
		Where("account_modified_event.name = ?", input.Name).
		Group("a.id", "a.account_id", "user_id")
	q := a.db.NewSelect().TableExpr("(?) as accounts", subQuery)
	if input.IsActiveOnly {
		q = q.
			Where("accounts.last_activated IS NOT NULL").
			WhereGroup(" AND ", func(query *bun.SelectQuery) *bun.SelectQuery {
				return query.Where("accounts.last_inactivated IS NULL").
					WhereOr("accounts.last_activated > accounts.last_inactivated")
			})
	}
	if err := q.Column("id", "account_id", "user_id").
		Order("last_activated DESC", "last_inactivated DESC").
		Limit(1).
		Scan(ctx, &result); err != nil {
		return nil, fmt.Errorf("failed to select account: %v", err)
	}
	account := &domains.Account{
		ID:        result.ID,
		AccountID: result.AccountID,
		UserID:    result.UserID,
	}
	return account, nil
}

func (a AccountService) Save(ctx context.Context, input domains.SaveAccountInput) error {
	a.logger.Infof("Save account start. input: %v", input)
	defer a.logger.Infof("Save account end. input: %v", input)
	// 1. 新規ならアカウント作成
	account := &domains.Account{
		AccountID: input.ID,
		UserID:    input.UserID,
	}
	accountExists, err := a.db.NewSelect().
		Model(account).
		Where("user_id = ?", input.UserID).
		Where("account_id = ?", input.ID).
		Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to select account: %v", err)
	}
	if !accountExists {
		if _, err := a.db.NewInsert().
			Model(account).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to update account: %v", err)
		}
	}
	if err := a.db.NewSelect().
		Model(account).
		Where("user_id = ?", input.UserID).
		Where("account_id = ?", input.ID).
		Limit(1).
		Scan(ctx); err != nil {
		return fmt.Errorf("failed to select account: %v", err)
	}

	// 2. アクティブならアクティブアカウントを登録、非アクティブなら非アクティブアカウントを登録
	if input.Active == 1 {
		if _, err := a.db.NewInsert().
			Model(&domains.ActiveAccount{AccountID: account.ID}).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert active account: %v", err)
		}
	} else if input.Active == -1 {
		if _, err := a.db.NewInsert().
			Model(&domains.InActiveAccount{AccountID: account.ID}).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert inactive account: %v", err)
		}
	}

	// 3. アカウント変更イベントを登録
	if _, err := a.db.NewInsert().
		Model(&domains.AccountModifiedEvent{
			AccountID: account.ID,
			Name:      input.Name,
			Sort:      input.Sort,
		}).
		Exec(ctx); err != nil {
		return fmt.Errorf("failed to insert account modified event: %v", err)
	}
	return nil
}

func NewAccountService(db bun.Tx, logger echo.Logger) *AccountService {
	return &AccountService{
		baseService: newBaseService(db, logger),
	}
}

var _ domains.AccountRepository = (*AccountService)(nil)
