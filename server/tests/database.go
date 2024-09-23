package tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/kanade0404/zaim/server/driver"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
	"path/filepath"
	"runtime"
	"testing"
)

func SetUp(ctx context.Context, t *testing.T) (*postgres.PostgresContainer, *bun.DB, error) {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		return nil, nil, errors.New("failed to get current file path")
	}
	basePath := filepath.Dir(b)
	f, err := filepath.Abs(filepath.Join(basePath, "..", "database", "migrations", "schema.sql"))
	if err != nil {
		return nil, nil, err
	}
	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:latest",
		postgres.WithInitScripts(f),
		postgres.WithDatabase("test_zaim"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to start postgres container: %v", err)
	}
	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Errorf("failed to terminate: %v", err)
		}
	})
	if err := postgresContainer.Snapshot(ctx, postgres.WithSnapshotName(t.Name())); err != nil {
		return postgresContainer, nil, err
	}
	dbURL, err := postgresContainer.ConnectionString(ctx)
	if err != nil {
		return nil, nil, err
	}
	db := driver.NewDB(dbURL + "sslmode=disable")
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	return postgresContainer, db, err
}
