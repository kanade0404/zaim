package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kanade0404/zaim/server/domains"
	"github.com/kanade0404/zaim/server/driver"
	"github.com/uptrace/bun"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	db := driver.NewDB(os.Getenv("DATABASE_URL"))
	if err := db.Ping(); err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	ctx := context.Background()
	if err := mode(ctx, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := user(ctx, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return nil
}

func mode(ctx context.Context, tx bun.Tx) error {
	if exists, err := tx.NewSelect().Model((*domains.Mode)(nil)).Where("id = ?", domains.PaymentMode).Exists(ctx); err != nil {
		return err
	} else {
		if !exists {
			if _, err := tx.NewInsert().Model(&domains.Mode{ID: domains.PaymentMode}).Exec(ctx); err != nil {
				return err
			}
		}
	}
	if exists, err := tx.NewSelect().Model((*domains.Mode)(nil)).Where("id = ?", domains.IncomeMode).Exists(ctx); err != nil {
		return err
	} else {
		if !exists {
			if _, err := tx.NewInsert().Model(&domains.Mode{ID: domains.IncomeMode}).Exec(ctx); err != nil {
				return err
			}
		}
	}
	if exists, err := tx.NewSelect().Model((*domains.Mode)(nil)).Where("id = ?", domains.TransferMode).Exists(ctx); err != nil {
		return err
	} else {
		if !exists {
			if _, err := tx.NewInsert().Model(&domains.Mode{ID: domains.TransferMode}).Exec(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

type UserSeedValue struct {
	ID                int64  `json:"id"`
	ConsumerKey       string `json:"consumer_key"`
	ConsumerSecret    string `json:"consumer_secret"`
	AccessToken       string `json:"oauth_token"`
	AccessTokenSecret string `json:"oauth_token_secret"`
}

type UserSeed map[string]*UserSeedValue

func user(ctx context.Context, tx bun.Tx) error {
	f, err := filepath.Abs("cmd/bun/seed/zaim.json")
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}
	b, err := os.ReadFile(f)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	var userSeed UserSeed
	decoder := json.NewDecoder(bytes.NewReader(b))
	if err := decoder.Decode(&userSeed); err != nil {
		return fmt.Errorf("failed to decode json: %v", err)
	}
	for k, seed := range userSeed {
		// user作成
		user := &domains.User{
			ID:   seed.ID,
			Name: k,
		}
		if _, err := tx.NewInsert().Model(user).On("CONFLICT (id) DO NOTHING").Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert user: %v", err)
		}
		if err := tx.NewSelect().Model(user).Where("id = ?", seed.ID).Limit(1).Scan(ctx); err != nil {
			return fmt.Errorf("failed to select user: %v", err)
		}
		if user == nil {
			return errors.New("failed to insert user")
		}
		// zaimアプリケーション作成
		zaimApp := &domains.ZaimApplication{
			UserID:         user.ID,
			ConsumerKey:    seed.ConsumerKey,
			ConsumerSecret: seed.ConsumerSecret,
		}
		if _, err := tx.NewInsert().Model(zaimApp).On("CONFLICT (user_id, consumer_key, consumer_secret) DO NOTHING").Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert zaim application: %v", err)
		}
		if err := tx.NewSelect().Model(zaimApp).Where("user_id = ?", user.ID).Where("consumer_key = ?", seed.ConsumerKey).Where("consumer_secret = ?", seed.ConsumerSecret).Limit(1).Scan(ctx); err != nil {
			return fmt.Errorf("failed to select zaim application: %v", err)
		}
		if zaimApp == nil {
			return errors.New("failed to insert zaim application")
		}
		// Zaimアプリケーション設定の有効化
		if exists, err := tx.NewSelect().Model((*domains.EnableZaimApplicationEvent)(nil)).Where("zaim_app_id = ?", zaimApp.ID).Exists(ctx); err != nil {
			return fmt.Errorf("failed to select enable zaim application event: %v", err)
		} else {
			if !exists {
				if _, err := tx.NewInsert().Model(&domains.EnableZaimApplicationEvent{
					ZaimAppID: zaimApp.ID,
				}).Exec(ctx); err != nil {
					return fmt.Errorf("failed to insert enable zaim application event: %v", err)
				}
			}
		}
		// zaim oauth認証作成
		zaimOauth := &domains.ZaimOAuth{
			ZaimAppID: zaimApp.ID,
			Token:     seed.AccessToken,
			Secret:    seed.AccessTokenSecret,
		}
		if _, err := tx.NewInsert().Model(zaimOauth).On("CONFLICT (zaim_app_id, token, secret) DO NOTHING").Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert zaim oauth: %v", err)
		}
		if err := tx.NewSelect().Model(zaimOauth).Where("zaim_app_id = ?", zaimApp.ID).Limit(1).Scan(ctx); err != nil {
			return fmt.Errorf("failed to select zaim oauth: %v", err)
		}
		if zaimOauth == nil {
			return errors.New("failed to insert zaim oauth")
		}
		// Config OAuth設定の有効化
		if exists, err := tx.NewSelect().Model((*domains.EnableZaimOAuthEvent)(nil)).Where("zaim_oauth_id = ?", zaimOauth.ID).Exists(ctx); err != nil {
			return fmt.Errorf("failed to select enable zaim oauth event: %v", err)
		} else {
			if !exists {
				if _, err := tx.NewInsert().Model(&domains.EnableZaimOAuthEvent{
					ZaimOAuthID: zaimOauth.ID,
				}).Exec(ctx); err != nil {
					return fmt.Errorf("failed to insert enable zaim application event: %v", err)
				}
			}
		}
	}
	return nil
}
