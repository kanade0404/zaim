package zaim

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"zaim/infrastructures/gcs"
	"zaim/infrastructures/zaim"
	"zaim/middlewares"
)

func ListActiveCategory(c echo.Context) error {
	ctx := c.(*middlewares.CustomContext)
	results := make(map[string][]zaim.Category)
	var errs []error
	for userName, config := range ctx.Config {
		zaimClient, err := zaim.NewClient(config.OAuthConfig.ConsumerKey, config.OAuthConfig.ConsumerSecret, config.OAuthToken.Token, config.OAuthToken.Secret)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		categories, err := zaimClient.ListActiveCategory()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		results[userName] = categories
	}
	if len(errs) > 0 {
		return c.JSON(http.StatusInternalServerError, errors.Join(errs...))
	}
	if err := gcs.PutCategory(c.Request().Context(), results); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, results)
}
