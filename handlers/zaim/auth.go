package zaim

import (
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"net/http"
	"net/url"
	"os"
	redis_client "zaim/infrastructures/redis"
	"zaim/middlewares"
)

func Authorization(c echo.Context) error {
	c.Logger().Info("start handler/authorization")
	defer c.Logger().Info("end handler/authorization")
	ctx := c.(*middlewares.CustomContext)
	results := make(map[string]string, len(ctx.Config))
	opt, err := redis.ParseURL(os.Getenv("REDIS_ENDPOINT"))
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	redisClient := redis.NewClient(opt)

	for key, cfg := range ctx.Config {
		fmt.Println(key, cfg)
		reqToken, reqSecret, err := cfg.OAuthConfig.RequestToken()
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		c.Logger().Debugf("reqToken: %s", reqToken)
		c.Logger().Debugf("reqSecret: %s", reqSecret)
		authorizeURL, err := cfg.OAuthConfig.AuthorizationURL(reqToken)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		authURL, err := url.Parse(authorizeURL.String())
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		if err := redis_client.SetZaimSecret(c.Request().Context(), redisClient, authURL.Query().Get("oauth_token"), key, reqSecret); err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		c.Logger().Debugf("authorizeURL: %s", authorizeURL)
		results[key] = authorizeURL.String()
	}

	return c.JSON(http.StatusOK, results)
}

func CallbackOAuthToken(c echo.Context) error {
	c.Logger().Info("start handler/callback")
	defer c.Logger().Info("end handler/callback")
	ctx := c.(*middlewares.CustomContext)
	opt, err := redis.ParseURL(os.Getenv("REDIS_ENDPOINT"))
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	redisClient := redis.NewClient(opt)
	requestToken, verifier, err := oauth1.ParseAuthorizationCallback(c.Request())
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Logger().Debugf("requestToken: %s", requestToken)
	c.Logger().Debugf("verifier: %s", verifier)
	results := make(map[string]string, len(ctx.Config))
	var hasSuccess bool
	oauthTokenParam := ctx.Request().URL.Query().Get("oauth_token")
	for key, cfg := range ctx.Config {
		requestSecret, err := redis_client.GetRequestSecret(c.Request().Context(), redisClient, oauthTokenParam)
		if err != nil {
			c.Logger().Error(err)
			results[key] = err.Error()
			continue
		}
		if key != requestSecret.User {
			continue
		}
		accessToken, accessSecret, err := cfg.OAuthConfig.AccessToken(requestToken, requestSecret.Secret, verifier)
		if err != nil {
			c.Logger().Error(err)
			results[key] = err.Error()
			continue
		}
		c.Logger().Debugf("accessToken: %s", accessToken)
		c.Logger().Debugf("accessSecret: %s", accessSecret)
		oauthToken := oauth1.NewToken(accessToken, accessSecret)
		c.Logger().Debugf("oauthToken: %s", oauthToken)
		if err := redis_client.SetOAuthTokens(ctx.Request().Context(), redisClient, key, oauthToken.Token, oauthToken.TokenSecret); err != nil {
			c.Logger().Error(err)
			results[key] = err.Error()
			continue
		}
		results[key] = "OK"
		hasSuccess = true
	}
	if !hasSuccess {
		return c.JSON(http.StatusInternalServerError, results)
	}
	defer func(ctx *middlewares.CustomContext, param string) {
		if err := redis_client.DeleteRequestSecret(ctx.Request().Context(), redisClient, param); err != nil {
			c.Logger().Error(err)
		}
	}(ctx, oauthTokenParam)
	return c.JSON(http.StatusOK, "OK")
}
