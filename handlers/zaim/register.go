package zaim

import (
	"encoding/csv"
	"fmt"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/drive/v3"
	"io"
)

func Register(c echo.Context) error {
	// csvをgoogle driveから取得する
	jstNow := synchro.Now[tz.AsiaTokyo]()
	// jstNowを先月の1日にする
	jstNow = jstNow.AddDate(0, -1, -jstNow.Day()+1)
	csvFileName := fmt.Sprintf("B43まとめ一覧_%d", jstNow.Year())
	srv, err := drive.NewService(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	r, err := srv.Files.List().Q(fmt.Sprintf("'1_bwbjpQK44ac5Eg9DGJ_Us3qE4eX1epq' in parents and name = '%s'", csvFileName)).Do()
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	if len(r.Files) == 0 {
		err := fmt.Errorf("not found %s", csvFileName)
		c.Logger().Error(err)
		return err
	}
	resp, err := srv.Files.Get(r.Files[0].Id).Download()
	if err != nil {
		c.Logger().Error(err)
		return err
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
		return err
	}
	for _, record := range records {
		fmt.Printf("%#v\n", record)
	}
	// csvをparseする
	// カスタムルールに基づいて変換する
	// 変換できないものはエラーにする
	// 変換したものをzaimに登録する
	return nil
}
