package zaim

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"zaim/infrastructures/zaim"
	"zaim/middlewares"
)

func ListActiveAccount(c echo.Context) error {
	ctx := c.(*middlewares.CustomContext)
	results := make(map[string][]zaim.Account)
	var errs []error
	for key, config := range ctx.Config {
		zaimClient, err := zaim.NewClient(config.OAuthConfig.ConsumerKey, config.OAuthConfig.ConsumerSecret, config.OAuthToken.Token, config.OAuthToken.Secret)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		accounts, err := zaimClient.ListActiveAccount()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		results[key] = accounts
	}
	if len(errs) > 0 {
		return c.JSON(http.StatusInternalServerError, errors.Join(errs...))
	}
	return c.JSON(http.StatusOK, results)
}
