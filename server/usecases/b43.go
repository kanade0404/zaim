package usecases

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/kanade0404/zaim/server/commands"
	"github.com/kanade0404/zaim/server/domains"
	"github.com/kanade0404/zaim/server/driver"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// DateIndex 日時
const DateIndex = 0

// ContentIndex 内容
const ContentIndex = 1

// MemoIndex メモ
const MemoIndex = 3

// AmountIndex 金額
const AmountIndex = 4

// CategoryIndex カテゴリ
const CategoryIndex = 5

// CanInclude まとめと予算に含める
const CanInclude = 6

type B43UseCase interface {
	SyncB43Transactions(ctx context.Context, userName string, startDate synchro.Time[tz.AsiaTokyo], dryRun bool) error
}

type b43Usecase struct {
	logger echo.Logger
	uow    commands.B43UnitOfWork
}

type B43Csv struct {
	Date     synchro.Time[tz.AsiaTokyo] `csv:"日時"`
	Content  string                     `csv:"内容"`
	Memo     string                     `csv:"メモ"`
	Amount   int                        `csv:"金額"`
	Category string                     `csv:"カテゴリ"`
	Mode     string
}

func newB43Csv(date, content, memo, amount, category string) (B43Csv, error) {
	var errs []error
	a, err := strconv.Atoi(amount)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to convert amount to int: %w", err))
	}
	d, err := synchro.Parse[tz.AsiaTokyo]("2006/01/02 15:04:05", date)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to parse date: %w", err))
	}
	if len(errs) != 0 {
		return B43Csv{}, fmt.Errorf("failed to create B43Csv: %w", errors.Join(errs...))
	}
	var mode string
	if a > 0 {
		mode = domains.IncomeMode
	} else if a < 0 {
		mode = domains.PaymentMode
	}
	return B43Csv{
		Date:     d,
		Content:  content,
		Memo:     memo,
		Amount:   a,
		Category: category,
		Mode:     mode,
	}, nil
}

// SyncB43Transactions はB43の取引を同期する
func (b b43Usecase) SyncB43Transactions(ctx context.Context, userName string, startDate synchro.Time[tz.AsiaTokyo], dryRun bool) error {
	b.logger.Infof("start SyncB43Transactions. userName: %s, startDate: %s, dryRun: %t", userName, startDate.String(), dryRun)
	defer b.logger.Infof("end SyncB43Transactions. userName: %s, startDate: %s, dryRun: %t", userName, startDate.String(), dryRun)
	// ローカルからcsvを取得する
	p, err := filepath.Abs(fmt.Sprintf("database/%s.csv", userName))
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	f, err := os.Open(p)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	// csvをパースする
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read csv: %w", err)
	}
	var zaimTransactions []B43Csv
	for _, record := range records[1:] {
		b43Csv, err := newB43Csv(record[DateIndex], record[ContentIndex], record[MemoIndex], record[AmountIndex], record[CategoryIndex])
		if err != nil {
			return fmt.Errorf("failed to create B43Csv: %w", err)
		}
		if record[CanInclude] != "1" {
			continue
		}
		if b43Csv.Date.Before(startDate) {
			continue
		}
		nextMonth := startDate.AddDate(0, 1, 0)
		if b43Csv.Date.After(nextMonth) {
			continue
		}
		zaimTransactions = append(zaimTransactions, b43Csv)
	}
	var (
		account    *domains.Account
		userID     int64
		userConfig *domains.Config
	)
	for _, transaction := range zaimTransactions {
		if err := b.uow.Do(ctx, func(ctx context.Context, uowRepositoryManager commands.B43RepositoryManager) error {
			if userID == 0 {
				u, err := uowRepositoryManager.UserRepositoryManager().FindByName(ctx, userName)
				if err != nil {
					return fmt.Errorf("failed to find u: %w", err)
				}
				userID = u.ID
			}
			if userConfig == nil {
				cfg, err := uowRepositoryManager.ConfigRepositoryManager().FindByUserID(ctx, userID)
				if err != nil {
					return fmt.Errorf("failed to find config: %w", err)
				}
				userConfig = cfg
			}
			zaimGenre, err := uowRepositoryManager.GenreRepository().FindB43Genre(ctx, domains.FindB43Genre{
				UserID: userID,
			})
			if err != nil {
				return fmt.Errorf("failed to find genre: %w", err)
			}
			zaimGenreID := zaimGenre.GenreID
			CategoryID := zaimGenre.Category.CategoryID
			if account == nil {
				account, err = uowRepositoryManager.AccountRepositoryManager().FindB43(ctx, domains.FindB43Input{UserID: userID})
				if err != nil {
					return fmt.Errorf("failed to find account: %w", err)
				}
			}
			var path string
			if transaction.Amount > 0 {
				path = "home/money/income"
			} else {
				path = "home/money/payment"
			}
			var comment string
			if transaction.Memo == "" {
				comment = "zaim-api"
			} else {
				comment = strings.Join([]string{transaction.Memo, "zaim-api"}, ",")
			}
			params := createPaymentBody(createPaymentInput{
				CategoryID:    CategoryID,
				GenreID:       zaimGenreID,
				Amount:        transaction.Amount,
				Date:          transaction.Date,
				FromAccountID: account.AccountID,
				Comment:       comment,
				Name:          transaction.Content,
				Place:         transaction.Content,
			})
			zaimDriver, err := driver.NewZaimDriver(userConfig.OAuthConfig.ConsumerKey, userConfig.OAuthConfig.ConsumerSecret, userConfig.OAuthToken.Token, userConfig.OAuthToken.Secret)
			if err != nil {
				return fmt.Errorf("failed to new zaim driver: %w", err)
			}
			if dryRun {
				b.logger.Info("Dry Run.")
				return nil
			}
			res, err := zaimDriver.Post(path, params)
			if err != nil {
				return fmt.Errorf("failed to post payment: %w", err)
			}
			if res.StatusCode >= 300 {
				return errors.New(fmt.Sprintf("failed to post payment: %d.", res.StatusCode))
			}
			return nil
		}); err != nil {
			return fmt.Errorf("failed to do transaction: %w", err)
		}
	}
	return nil
}

func NewB43Usecase(db *bun.DB, logger echo.Logger) B43UseCase {
	return &b43Usecase{logger: logger, uow: commands.NewB43UnitOfWork(db, logger)}
}

// SharedFilter SharedTransactionFilter は共有する取引をフィルタリングする
func SharedFilter(record []string) bool {
	switch record[CategoryIndex] {
	case "ショッピング", "交通", "教育":
		return false
	}
	return true
}

func PrivateFilter(record []string) bool {
	switch record[ContentIndex] {
	case "AMAZON.CO.JP":
		return false
	}
	return true
}

type createPaymentInput struct {
	CategoryID    int64
	GenreID       int64
	Amount        int
	Date          synchro.Time[tz.AsiaTokyo]
	FromAccountID int64
	Comment       string
	Name          string
	Place         string
}

type CreatePaymentParams struct {
	CategoryID    int    `json:"category_id"`
	GenreID       int    `json:"genre_id"`
	Amount        int    `json:"amount"`
	Date          string `json:"date"`
	FromAccountID int    `json:"from_account_id"`
	Comment       string `json:"comment"`
	Name          string `json:"name"`
	Place         string `json:"place"`
}

func createPaymentBody(input createPaymentInput) CreatePaymentParams {
	var amount int
	if input.Amount >= 0 {
		amount = input.Amount
	} else {
		amount = -input.Amount
	}
	values := CreatePaymentParams{}
	values.CategoryID = int(input.CategoryID)
	values.GenreID = int(input.GenreID)
	values.Amount = amount
	values.Date = input.Date.Format("2006-01-02")
	if input.FromAccountID != 0 {
		values.FromAccountID = int(input.FromAccountID)
	}
	if input.Comment != "" {
		values.Comment = input.Comment
	}
	if input.Name != "" {
		values.Name = input.Name
	}
	if input.Place != "" {
		values.Place = input.Place
	}
	return values
}
