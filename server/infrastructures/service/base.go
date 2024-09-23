package service

import (
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type baseService struct {
	db     bun.IDB
	logger echo.Logger
}

func newBaseService(db bun.IDB, logger echo.Logger) baseService {
	return baseService{db: db, logger: logger}
}
func newTestBaseService(db bun.IDB) baseService {
	return baseService{db: db, logger: echo.New().Logger}
}
