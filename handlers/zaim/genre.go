package zaim

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"zaim/infrastructures/gcs"
	"zaim/infrastructures/zaim"
	"zaim/middlewares"
)

func ListActiveGenre(c echo.Context) error {
	ctx := c.(*middlewares.CustomContext)
	results := make(map[string][]zaim.Genre)
	var errs []error
	for key, config := range ctx.Config {
		zaimClient, err := zaim.NewClient(config.OAuthConfig.ConsumerKey, config.OAuthConfig.ConsumerSecret, config.OAuthToken.Token, config.OAuthToken.Secret)
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
	if len(errs) > 0 {
		return c.JSON(http.StatusInternalServerError, errs)
	}
	if err := gcs.PutGenre(c.Request().Context(), results); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, results)
}
