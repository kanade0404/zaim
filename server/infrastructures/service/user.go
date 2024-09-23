package service

import (
	"context"
	"github.com/kanade0404/zaim/server/domains"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type UserService struct {
	baseService
}

func (u UserService) FindByName(ctx context.Context, name string) (*domains.User, error) {
	user := new(domains.User)
	err := u.db.NewSelect().Model(user).Where("name = ?", name).Limit(1).Scan(ctx)
	return user, err
}

func NewUserService(db bun.Tx, logger echo.Logger) domains.UserRepository {
	return UserService{baseService: newBaseService(db, logger)}
}
