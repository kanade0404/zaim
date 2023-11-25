package zaim

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
	"zaim/handlers"
	"zaim/infrastructures/zaim"
	"zaim/usecases"
)

type body struct {
	RunAt  string `json:"run_at" form:"run_at"`
	DryRun bool   `json:"dry_run" form:"dry_run"`
}

type RegisterResponse struct {
	Responses [][]zaim.PaymentResponse `json:"responses"`
}

func Register(c echo.Context) error {
	var body body
	if err := c.Bind(&body); err != nil {
		return c.JSON(400, err)
	}
	runAt, err := synchro.Parse[tz.AsiaTokyo](time.DateOnly, body.RunAt)
	if err != nil {
		return err
	}
	// jstNowを先月の1日にする
	jstLastMonth := runAt.AddDate(0, -1, -runAt.Day()+1)
	res, err := usecases.RegisterMonthlyTransactions(c, jstLastMonth.StdTime(), body.DryRun)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Error: err,
		})
	}
	return c.JSON(http.StatusOK, RegisterResponse{Responses: res})
}
