package zaim

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/kanade0404/zaim/server/driver"
	"github.com/kanade0404/zaim/server/usecases"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"time"
)

func B43Register(c echo.Context) error {
	var (
		data struct {
			RunAt  *string `json:"run_at"`
			DryRun bool    `json:"dry_run"`
			User   string  `json:"user"`
		}
		runAt synchro.Time[tz.AsiaTokyo]
	)
	if err := c.Bind(&data); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.Logger().Error(err)
		}
	}(c.Request().Body)
	if data.RunAt == nil {
		runAt = synchro.Now[tz.AsiaTokyo]()
	} else {
		var err error
		runAt, err = synchro.Parse[tz.AsiaTokyo](time.DateOnly, *(data.RunAt))
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, err)
		}
	}
	ctx := c.Request().Context()
	db := driver.NewDB(os.Getenv("DATABASE_URL"))
	if err := db.Ping(); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err := usecases.NewB43Usecase(db, c.Logger()).SyncB43Transactions(ctx, data.User, runAt, data.DryRun); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "ok")
}
