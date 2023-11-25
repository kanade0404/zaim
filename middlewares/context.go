package middlewares

import (
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"os"
	"strings"
)

type Redis struct {
	*redis.Client
	Config map[string]*oauth1.Config
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
		users := strings.Split(os.Getenv("USERS"), ",")
		configs := make(map[string]*oauth1.Config, len(users))
		for _, user := range users {
			configs[user] = &oauth1.Config{
				ConsumerKey:    os.Getenv(fmt.Sprintf("%s_CONSUMER_KEY", strings.ToUpper(user))),
				ConsumerSecret: os.Getenv(fmt.Sprintf("%s_CONSUMER_SECRET", strings.ToUpper(user))),
				CallbackURL:    fmt.Sprintf("%s/auth/callback", os.Getenv("HOST")),
				Endpoint: oauth1.Endpoint{
					RequestTokenURL: requestURL,
					AuthorizeURL:    authorizeURL,
					AccessTokenURL:  accessURL,
				},
			}
		}
		c.Logger().Infof("configs: %v", configs)
		return next(&CustomContext{c, Redis{
			Client: redis.NewClient(opt),
			Config: configs,
		}})
	}
}
