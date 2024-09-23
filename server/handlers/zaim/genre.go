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

func UpdateGenres(c echo.Context) error {
	var b struct {
		User string `json:"user"`
	}
	if err := c.Bind(&b); err != nil {
		c.Logger().Errorf("failed to bind body: %v", err)
		return c.JSON(http.StatusBadRequest, fmt.Errorf("failed to bind body: %v", err))
	}
	ctx := c.Request().Context()
	db := driver.NewDB(os.Getenv("DATABASE_URL"))
	if err := db.Ping(); err != nil {
		c.Logger().Errorf("failed to ping db: %v", err)
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("failed to ping db: %v", err))
	}
	if err := usecases.NewUpdateGenreUseCase(db, c.Logger()).UpdateGenres(ctx, strings.ToUpper(b.User)); err != nil {
		c.Logger().Errorf("failed to update genre: %v", err)
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("failed to update genre: %v", err))
	}
	return c.JSON(http.StatusCreated, "success")
}
