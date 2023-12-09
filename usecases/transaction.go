package usecases

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/drive/v3"
	"io"
	"strconv"
	"strings"
	"time"
	"zaim/infrastructures/gcs"
	"zaim/infrastructures/zaim"
	"zaim/middlewares"
)

const DateIndex = 0
const ContentIndex = 1
const MemoIndex = 3
const AmountIndex = 4
const CategoryIndex = 5

const Spending = "消費"
const Investing = "投資"
const Wasting = "浪費"

func RegisterMonthlyTransactions(c echo.Context, jstLastMonth time.Time, isDryRun bool) ([][]zaim.PaymentResponse, error) {
	csvFileName := fmt.Sprintf("B43まとめ一覧_%d", jstLastMonth.Year())
	srv, err := drive.NewService(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return nil, err
	}
	ctx := c.(*middlewares.CustomContext)
	var (
		responses [][]zaim.PaymentResponse
		errs      []error
	)
	for k := range ctx.Config {
		responses, err := RegisterTransaction(c, srv, ctx.Config[k].CsvFolder, csvFileName, k, jstLastMonth, isDryRun)
		if err != nil {
			errs = append(errs, err)
		} else {
			responses = append(responses, responses...)
		}
	}
	return responses, errors.Join(errs...)
}

func RegisterTransaction(c echo.Context, srv *drive.Service, csvFolderName, csvFileName string, userName string, jstLastMonth time.Time, isDryRun bool) ([]zaim.PaymentResponse, error) {
	c.Logger().Infof("start usecases/register. userName: %s, jstLastMonth: %s, isDryRun: %t", userName, jstLastMonth.String(), isDryRun)
	defer c.Logger().Infof("end usecases/register. userName: %s, jstLastMonth: %s, isDryRun: %t", userName, jstLastMonth.String(), isDryRun)
	r, err := srv.Files.List().Q(fmt.Sprintf("'%s' in parents and name = '%s'", csvFolderName, csvFileName)).Do()
	if err != nil {
		c.Logger().Error(err)
		return nil, err
	}
	if len(r.Files) == 0 {
		err := fmt.Errorf("not found %s", csvFileName)
		c.Logger().Error(err)
		return nil, err
	}
	resp, err := srv.Files.Get(r.Files[0].Id).Download()
	if err != nil {
		c.Logger().Error(err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.Logger().Error(err)
		}
	}(resp.Body)
	obj := csv.NewReader(resp.Body)
	records, err := obj.ReadAll()
	if err != nil {
		c.Logger().Error(err)
		return nil, err
	}
	var payments []zaim.PaymentParameter
	// categoryをjsonから取得する
	categories, err := gcs.GetCategoryByUserName(c.Request().Context(), userName)
	if err != nil {
		c.Logger().Error(err)
		return nil, err
	}
	categoryMap := make(map[int]zaim.Category)
	for _, category := range categories {
		categoryMap[category.ID] = category
	}
	// genreをjsonから取得する
	genres, err := gcs.GetGenreByUserName(c.Request().Context(), userName)
	if err != nil {
		c.Logger().Error(err)
		return nil, err
	}
	genreMap := make(map[int]zaim.Genre)
	for _, genre := range genres {
		genreMap[genre.ID] = genre
	}
	// genre_mappingをjsonから取得する
	genreMappings, err := gcs.GetGenreMappingByUserName(c.Request().Context(), userName)
	if err != nil {
		c.Logger().Error(err)
		return nil, err
	}
	contentMappings, err := gcs.GetContentMappingByUserName(c.Request().Context(), userName)
	if userName == "SHARED" && err != nil {
		c.Logger().Error(err)
		return nil, err
	}
	var currentMonthRecords [][]string
	for _, record := range records[1:] {
		d, err := synchro.Parse[tz.AsiaTokyo]("2006/01/02 15:04:05", record[DateIndex])
		if err != nil {
			c.Logger().Errorf("failed to parse date '%s': %v, record: %v", record[DateIndex], err, record)
			continue
		}
		if d.Year() != jstLastMonth.Year() || d.Month() != jstLastMonth.Month() {
			continue
		}
		currentMonthRecords = append(currentMonthRecords, record)
	}
	// メモに'支払いキャンセルによる返金'とあるものは同じ名前で同じ金額のものを除外する
	var canceledPaymentRecords [][]string
	for _, record := range currentMonthRecords {
		if strings.Contains(record[MemoIndex], "支払いキャンセルによる返金") {
			canceledPaymentRecords = append(canceledPaymentRecords, record)
		}
	}
	canceledPaymentFoundCount := make(map[int]int)
	for canceledID, canceledPaymentRecord := range canceledPaymentRecords {
		for i, record := range currentMonthRecords {
			if canceledPaymentRecord[AmountIndex] == strings.Replace(record[AmountIndex], "-", "", 1) && canceledPaymentRecord[ContentIndex] == record[ContentIndex] {
				if _, ok := canceledPaymentFoundCount[i]; !ok {
					canceledPaymentFoundCount[i] = canceledID
				}
				currentMonthRecords = append(currentMonthRecords[:i], currentMonthRecords[i+1:]...)
				break
			}
		}
	}
	if len(canceledPaymentFoundCount) != len(canceledPaymentRecords) {
		err := fmt.Errorf("canceled payment count is mismatch. canceledPaymentFoundCount: %d, canceledPaymentRecords: %d", canceledPaymentFoundCount, len(canceledPaymentRecords))
		c.Logger().Error(err)
		return nil, err

	}
	// 日時,内容,操作者,メモ,金額,カテゴリ,まとめと予算に含める
	for _, record := range records[1:] {
		d, err := synchro.Parse[tz.AsiaTokyo]("2006/01/02 15:04:05", record[DateIndex])
		if err != nil {
			c.Logger().Errorf("failed to parse date '%s': %v, record: %v", record[DateIndex], err, record)
			continue
		}
		if d.Year() != jstLastMonth.Year() || d.Month() != jstLastMonth.Month() {
			continue
		}
		if userName == "SHARED" {
			if !SharedFilter(record) {
				continue
			}
		} else {
			if !PrivateFilter(record) {
				continue
			}
		}
		b43Category := record[CategoryIndex]
		datetime, err := synchro.Parse[tz.AsiaTokyo]("2006/01/02 15:04:05", record[DateIndex])
		if err != nil {
			c.Logger().Errorf("failed to parse date '%s': %v, record: %v", record[DateIndex], err, record)
			continue
		}
		date := datetime.Format(time.DateOnly)
		amount := strings.Replace(record[AmountIndex], "-", "", 1)
		content := record[ContentIndex]
		memo := record[MemoIndex]
		if m, ok := genreMappings[b43Category]; ok {
			if genre, ok := genreMap[m.GenreID]; ok {
				if category, ok := categoryMap[genre.CategoryID]; ok {
					payments = append(payments, zaim.PaymentParameter{
						Date:       date,
						CategoryID: strconv.Itoa(category.ID),
						GenreID:    strconv.Itoa(genre.ID),
						Amount:     amount,
						Name:       content,
						Place:      content,
						Comment:    memo,
					})
				} else {
					c.Logger().Infof("not found category: %v", record)
					continue
				}
			} else {
				c.Logger().Infof("not found genre: %v", record)
				continue
			}
		} else {
			if strings.Contains(record[ContentIndex], "Spotify") {
				payments = append(payments, zaim.PaymentParameter{
					Date:       date,
					CategoryID: "55982027",
					GenreID:    "32551099",
					Amount:     amount,
					Name:       content,
					Place:      content,
					Comment:    memo,
				})
			} else if m, ok := contentMappings[record[ContentIndex]]; ok {
				payments = append(payments, zaim.PaymentParameter{
					Date:       date,
					CategoryID: strconv.Itoa(m.CategoryID),
					GenreID:    strconv.Itoa(m.GenreID),
					Amount:     amount,
					Name:       content,
					Place:      content,
					Comment:    memo,
				})
			} else if strings.Contains(b43Category, Spending) {
				payments = append(payments, zaim.PaymentParameter{
					Date:       date,
					CategoryID: strconv.Itoa(genreMappings["その他消費"].CategoryID),
					GenreID:    strconv.Itoa(genreMappings["その他消費"].GenreID),
					Amount:     amount,
					Name:       content,
					Place:      content,
					Comment:    memo,
				})
			} else if strings.Contains(b43Category, Investing) {
				payments = append(payments, zaim.PaymentParameter{
					Date:       date,
					CategoryID: strconv.Itoa(genreMappings["その他投資"].CategoryID),
					GenreID:    strconv.Itoa(genreMappings["その他投資"].GenreID),
					Amount:     amount,
					Name:       content,
					Place:      content,
					Comment:    memo,
				})
			} else if strings.Contains(b43Category, Wasting) {
				payments = append(payments, zaim.PaymentParameter{
					Date:       date,
					CategoryID: strconv.Itoa(genreMappings["その他浪費"].CategoryID),
					GenreID:    strconv.Itoa(genreMappings["その他浪費"].GenreID),
					Amount:     amount,
					Name:       content,
					Place:      content,
					Comment:    memo,
				})
			} else {
				c.Logger().Error("not found genre_mapping: %v", record)
				continue
			}
		}
	}
	ctx := c.(*middlewares.CustomContext)
	oauthToken := ctx.Config[userName].OAuthToken.Token
	oauthTokenSecret := ctx.Config[userName].OAuthToken.Secret
	zaimClient, err := zaim.NewClient(ctx.Config[userName].OAuthConfig.ConsumerKey, ctx.Config[userName].OAuthConfig.ConsumerSecret, oauthToken, oauthTokenSecret)
	if err != nil {
		c.Logger().Error(err)
		return nil, err
	}
	var responses []zaim.PaymentResponse
	var paymentErrs []error
	for _, payment := range payments {
		payment.Comment = strings.Join([]string{payment.Comment, "zaim-api"}, ",")
		if userName == "SHARED" {
			payment.FromAccountID = "17647716"
		} else {
			payment.FromAccountID = "17326014"
		}
		if isDryRun {
			c.Logger().Infof("Dry Run. payment: %#v\n", payment)
			fmt.Printf("Dry Run. payment: %#v\n", payment)
			continue
		}
		res, err := zaimClient.CreatePayment(payment)
		if err != nil {
			c.Logger().Error(err)
			paymentErrs = append(paymentErrs, err)
			continue
		}
		responses = append(responses, res)
	}
	return responses, errors.Join(paymentErrs...)
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
