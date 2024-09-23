package driver

import (
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func NewDB(dbURL string) *bun.DB {
	return bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbURL))), pgdialect.New())
}

func Rollback(tx bun.Tx, err error) error {
	if err := tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback: %v\n", err)
	}
	return err

}
