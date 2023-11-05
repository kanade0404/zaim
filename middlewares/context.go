package middlewares

import (
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"os"
)

type Redis struct {
	*redis.Client
	Config *oauth1.Config
}
type CustomContext struct {
	echo.Context
	Redis Redis
}

const providerBaseURL = "https://api.zaim.net/v2/auth/%s"

var requestURL = fmt.Sprintf(providerBaseURL, "request")

const authorizeURL = "https://auth.zaim.net/users/auth"

var accessURL = fmt.Sprintf(providerBaseURL, "access")

func Context(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().Info("start middleware/context")
		defer c.Logger().Info("end middleware/context")
		c.Logger().Debug("Custom middleware")
		opt, err := redis.ParseURL(os.Getenv("REDIS_ENDPOINT"))
		if err != nil {
			c.Logger().Fatal(err)
		}
		return next(&CustomContext{c, Redis{
			Client: redis.NewClient(opt),
			Config: &oauth1.Config{
				ConsumerKey:    os.Getenv("CONSUMER_KEY"),
				ConsumerSecret: os.Getenv("CONSUMER_SECRET"),
				CallbackURL:    fmt.Sprintf("%s/callback", os.Getenv("HOST")),
				Endpoint: oauth1.Endpoint{
					RequestTokenURL: requestURL,
					AuthorizeURL:    authorizeURL,
					AccessTokenURL:  accessURL,
				},
			},
		}})
	}
}
