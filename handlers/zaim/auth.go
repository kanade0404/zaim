package zaim

import (
	"github.com/dghubble/oauth1"
	"github.com/labstack/echo/v4"
	"net/http"
	"zaim/infrastructures/redis"
	"zaim/middlewares"
)

func Authorization(c echo.Context) error {
	c.Logger().Info("start handler/authorization")
	defer c.Logger().Info("end handler/authorization")
	ctx := c.(*middlewares.CustomContext)
	reqToken, reqSecret, err := ctx.Redis.Config.RequestToken()
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Logger().Debugf("reqToken: %s", reqToken)
	c.Logger().Debugf("reqSecret: %s", reqSecret)
	if err := redis.SetRequestSecret(c.Request().Context(), ctx.Redis.Client, reqSecret); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	authorizeURL, err := ctx.Redis.Config.AuthorizationURL(reqToken)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Logger().Debugf("authorizeURL: %s", authorizeURL)
	return c.JSON(http.StatusOK, struct {
		URL string `json:"url"`
	}{
		authorizeURL.String(),
	})
}

func CallbackOAuthToken(c echo.Context) error {
	c.Logger().Info("start handler/callback")
	defer c.Logger().Info("end handler/callback")
	ctx := c.(*middlewares.CustomContext)
	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(c.Request())
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Logger().Debugf("requestToken: %s", requestToken)
	c.Logger().Debugf("verifier: %s", verifier)
	requestSecret, err := redis.GetRequestSecret(c.Request().Context(), ctx.Redis.Client)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	accessToken, accessSecret, err := ctx.Redis.Config.AccessToken(requestToken, requestSecret, verifier)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Logger().Debugf("accessToken: %s", accessToken)
	c.Logger().Debugf("accessSecret: %s", accessSecret)
	oauthToken := oauth1.NewToken(accessToken, accessSecret)
	c.Logger().Debugf("oauthToken: %s", oauthToken)
	if err := redis.SetOAuthToken(c.Request().Context(), ctx.Redis.Client, oauthToken.Token); err != nil {
		c.Logger().Error(err)
		return err
	}
	if err := redis.SetOAuthTokenSecret(c.Request().Context(), ctx.Redis.Client, oauthToken.TokenSecret); err != nil {
		c.Logger().Error(err)
		return err
	}
	return c.JSON(http.StatusOK, "OK")
}
