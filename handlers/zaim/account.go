package zaim

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"zaim/infrastructures/redis"
	"zaim/infrastructures/zaim"
	"zaim/middlewares"
)

func ListActiveAccount(c echo.Context) error {
	ctx := c.(*middlewares.CustomContext)
	configs := ctx.Redis.Config
	results := make(map[string][]zaim.Account)
	var (
		errs []error
	)
	for key, config := range configs {
		oauthToken, err := redis.GetOauthToken(ctx.Request().Context(), ctx.Redis.Client, key)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		zaimClient, err := zaim.NewClient(config.ConsumerKey, config.ConsumerSecret, oauthToken.Token, oauthToken.Secret)
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
		fmt.Println(errs)
		return c.JSON(http.StatusInternalServerError, errors.Join(errs...))
	}
	return c.JSON(http.StatusOK, results)
}
