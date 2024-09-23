package main

import (
	"fmt"
	"github.com/kanade0404/zaim/server/domains"
	"github.com/kanade0404/zaim/server/driver"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
	"log"
	"os"
	"path/filepath"
)

var models = []domains.Domain{
	(*domains.User)(nil),
	(*domains.ZaimApplication)(nil),
	(*domains.ZaimOAuth)(nil),
	(*domains.EnableZaimApplicationEvent)(nil),
	(*domains.EnableZaimOAuthEvent)(nil),

	(*domains.Account)(nil),
	(*domains.AccountModifiedEvent)(nil),
	(*domains.ActiveAccount)(nil),
	(*domains.InActiveAccount)(nil),
	(*domains.B43Account)(nil),

	(*domains.Mode)(nil),

	(*domains.Category)(nil),
	(*domains.CategoryModifiedEvent)(nil),
	(*domains.ActiveCategory)(nil),
	(*domains.InActiveCategory)(nil),

	(*domains.Genre)(nil),
	(*domains.GenreModifiedEvent)(nil),
	(*domains.ActiveGenre)(nil),
	(*domains.InActiveGenre)(nil),
	(*domains.ParentGenreModifiedEvent)(nil),
	(*domains.B43Genre)(nil),
}

func generateUpQuery(db *bun.DB, models []domains.Domain) []byte {
	var data []byte
	for _, model := range models {
		query := db.NewCreateTable().Model(model).WithForeignKeys()
		rawQuery, err := query.AppendQuery(db.Formatter(), nil)
		if err != nil {
			panic(err)
		}
		data = append(data, rawQuery...)
		data = append(data, ";\n"...)
	}
	for _, model := range models {
		for _, index := range model.Indexes() {
			idx := db.NewCreateIndex().Model(model).Index(index.Name()).Column(index.Columns()...)
			rawQuery, err := idx.AppendQuery(db.Formatter(), nil)
			if err != nil {
				panic(err)
			}
			data = append(data, rawQuery...)
			data = append(data, ";\n"...)
		}
	}
	return data
}

func generateDownQuery(db *bun.DB, models []domains.Domain) []byte {
	var data []byte
	for i := len(models) - 1; i >= 0; i-- {
		for _, index := range models[i].Indexes() {
			idx := db.NewDropIndex().Model(models[i]).Index(index.Name())
			rawQuery, err := idx.AppendQuery(db.Formatter(), nil)
			if err != nil {
				panic(err)
			}
			data = append(data, rawQuery...)
			data = append(data, ";\n"...)
		}
	}
	for i := len(models) - 1; i >= 0; i-- {
		query := db.NewDropTable().Model(models[i])
		rawQuery, err := query.AppendQuery(db.Formatter(), nil)
		if err != nil {
			panic(err)
		}
		data = append(data, rawQuery...)
		data = append(data, ";\n"...)
	}
	return data
}

func main() {
	if err := migration(); err != nil {
		log.Fatal(err)
	}
}

func migration() error {
	db := driver.NewDB(os.Getenv("DATABASE_URL"))
	if os.Getenv("ENV") == "local" {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}
	if err := db.Ping(); err != nil {
		return err
	}
	p, err := filepath.Abs("./database/migrations")
	if err != nil {
		return err
	}
	if err := os.WriteFile(fmt.Sprintf("%s/schema.sql", p), generateUpQuery(db, models), 0666); err != nil {
		return err
	}
	if err := os.WriteFile(fmt.Sprintf("%s/down.sql", p), generateDownQuery(db, models), 0666); err != nil {
		return err
	}
	return nil
}
