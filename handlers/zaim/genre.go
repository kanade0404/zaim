package zaim

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"zaim/infrastructures/redis"
	"zaim/infrastructures/zaim"
	"zaim/middlewares"
)

func ListActiveGenre(c echo.Context) error {
	ctx := c.(*middlewares.CustomContext)
	configs := ctx.Redis.Config
	results := make(map[string][]zaim.Genre)
	var errs []error
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
		genres, err := zaimClient.ListActiveGenre()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		results[key] = genres
	}
	return c.JSON(http.StatusOK, results)
}
