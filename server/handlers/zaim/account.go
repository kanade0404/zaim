package zaim

import (
	"fmt"
	"github.com/kanade0404/zaim/server/driver"
	"github.com/kanade0404/zaim/server/usecases"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
)

func UpdateAccount(c echo.Context) error {
	var b struct {
		User string `json:"user"`
	}
	if err := c.Bind(&b); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("failed to bind body: %v", err))
	}
	zaimUserName := strings.ToUpper(b.User)
	ctx := c.Request().Context()
	db := driver.NewDB(os.Getenv("DATABASE_URL"))
	if err := db.Ping(); err != nil {
		c.Logger().Errorf("failed to ping db: %v", err)
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("failed to ping db: %v", err))
	}

	if err := usecases.NewUpdateAccountUseCase(db, c.Logger()).UpdateAccounts(ctx, zaimUserName); err != nil {
		c.Logger().Errorf("failed to update accounts: %v", err)
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("failed to update accounts: %v", err))
	}
	return c.JSON(http.StatusCreated, nil)
}
