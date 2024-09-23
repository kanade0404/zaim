package domains

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Account struct {
	bun.BaseModel         `bun:"table:account,alias:a"`
	ID                    uuid.UUID               `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	AccountID             int64                   `bun:"account_id,notnull,unique:unq_account"`
	UserID                int64                   `bun:"user_id,notnull,unique:unq_account"`
	User                  *User                   `bun:"rel:belongs-to,join:user_id=id"`
	AccountModifiedEvents []*AccountModifiedEvent `bun:"rel:has-many,join:id=account_id"`
	ActiveAccounts        []*ActiveAccount        `bun:"rel:has-many,join:id=account_id"`
	InActiveAccounts      []*InActiveAccount      `bun:"rel:has-many,join:id=account_id"`
	ZaimAccounts          []*B43Account           `bun:"rel:has-many,join:id=account_id"`
}

var _ Domain = (*Account)(nil)

func (a *Account) Indexes() []*Index {
	return nil
}

type AccountModifiedEvent struct {
	bun.BaseModel `bun:"table:account_modified_event,alias:ame"`
	ID            uuid.UUID `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	AccountID     uuid.UUID `bun:"account_id,type:uuid,notnull,unique:unq_account_modified_event"`
	Account       *Account  `bun:"rel:belongs-to,join:account_id=id"`
	Name          string    `bun:"name,notnull,unique:unq_account_modified_event"`
	Sort          int       `bun:"sort,notnull,unique:unq_account_modified_event"`
	Modified      time.Time `bun:"modified,notnull,default:current_timestamp,unique:unq_account_modified_event"`
}

var _ Domain = (*AccountModifiedEvent)(nil)

func (m *AccountModifiedEvent) Indexes() []*Index {
	return []*Index{
		createIndex("account_modified_event", []string{"account_id"}),
		createIndex("account_modified_event", []string{"name"}),
	}
}

type ActiveAccount struct {
	bun.BaseModel `bun:"table:active_account,alias:aa"`
	ID            uuid.UUID `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	AccountID     uuid.UUID `bun:"account_id,type:uuid,notnull,unique:unq_active_account"`
	Account       *Account  `bun:"rel:belongs-to,join:account_id=id"`
	Activated     time.Time `bun:"activated,notnull,default:current_timestamp,unique:unq_active_account"`
}

var _ Domain = (*ActiveAccount)(nil)

func (a *ActiveAccount) Indexes() []*Index {
	return nil
}

type InActiveAccount struct {
	bun.BaseModel `bun:"table:inactive_account,alias:ia"`
	ID            uuid.UUID `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	AccountID     uuid.UUID `bun:"account_id,type:uuid,notnull,unique:unq_inactive_account"`
	Account       *Account  `bun:"rel:belongs-to,join:account_id=id"`
	InActivated   time.Time `bun:"inactivated,notnull,default:current_timestamp,unique:unq_inactive_account"`
}

var _ Domain = (*InActiveAccount)(nil)

func (a *InActiveAccount) Indexes() []*Index {
	return nil
}

type B43Account struct {
	bun.BaseModel `bun:"table:b43_account,alias:ba"`
	ID            uuid.UUID `bun:"id,type:uuid,pk,default:gen_random_uuid()"`
	AccountID     uuid.UUID `bun:"account_id,type:uuid,notnull,unique:unq_b43_account"`
	Account       *Account  `bun:"rel:belongs-to,join:account_id=id"`
	Modified      time.Time `bun:"modified,notnull,default:current_timestamp,unique:unq_b43_account"`
}

func (z *B43Account) Indexes() []*Index {
	return nil
}

var _ Domain = (*B43Account)(nil)

type FindByNameAccountInput struct {
	Name         string
	UserID       int64
	IsActiveOnly bool
}
type SaveAccountInput struct {
	ID     int64
	Name   string
	Sort   int
	Active int
	UserID int64
}
type FindB43Input struct {
	UserID          int64
	IncludeInActive bool
}

type AccountRepository interface {
	Save(ctx context.Context, input SaveAccountInput) error
	FindByName(ctx context.Context, input FindByNameAccountInput) (*Account, error)
	FindB43(ctx context.Context, input FindB43Input) (*Account, error)
}
