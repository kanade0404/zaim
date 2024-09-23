package domains

import (
	"context"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:user,alias:u"`
	ID            int64              `bun:"id,pk"`
	Name          string             `bun:"name,unique,notnull"`
	ZaimApps      []*ZaimApplication `bun:"rel:has-many,join:id=user_id"`
	Categories    []*Category        `bun:"rel:has-many,join:id=user_id"`
	Accounts      []*Account         `bun:"rel:has-many,join:id=user_id"`
}

func (u *User) Indexes() []*Index {
	return []*Index{
		createIndex("user", []string{"name"}),
	}
}

type UserRepository interface {
	FindByName(ctx context.Context, name string) (*User, error)
}
