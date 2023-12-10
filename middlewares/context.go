package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/labstack/echo/v4"
	"io"
	"os"
	"zaim/infrastructures/secret_manager"
)

type OAuthToken struct {
	Token  string
	Secret string
}
type Zaim struct {
	OAuthConfig *oauth1.Config
	OAuthToken
	CsvFolder string
}
type CustomContext struct {
	echo.Context
	Config map[string]Zaim
}

const providerBaseURL = "https://api.zaim.net/v2/auth/%s"

var requestURL = fmt.Sprintf(providerBaseURL, "request")

const authorizeURL = "https://auth.zaim.net/users/auth"

var accessURL = fmt.Sprintf(providerBaseURL, "access")

type body struct {
	Users []string `json:"users"`
}

func Context(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().Info("start middleware/context")
		defer c.Logger().Info("end middleware/context")
		// 後続でstreamを消費しないようにする
		req := c.Request()
		ctx := req.Context()
		var body body
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				c.Logger().Error(err)
			}
		}(req.Body)
		b, err := io.ReadAll(req.Body)
		if err != nil {
			c.Logger().Fatal(err)
		}
		if err := json.Unmarshal(b, &body); err != nil {
			c.Logger().Fatal(err)
		}
		users := body.Users
		secretDriver, err := secret_manager.NewDriver(ctx)
		if err != nil {
			c.Logger().Fatal(err)
		}
		configs := make(map[string]Zaim, len(users))
		for _, user := range users {
			consumerKey, err := secretDriver.GetConsumerKey(user)
			if err != nil {
				c.Logger().Fatal(err)
			}
			consumerSecret, err := secretDriver.GetConsumerSecret(user)
			if err != nil {
				c.Logger().Fatal(err)
			}
			oauthToken, err := secretDriver.GetOAuthToken(user)
			if err != nil {
				c.Logger().Fatal(err)
			}
			oauthSecret, err := secretDriver.GetOAuthSecret(user)
			if err != nil {
				c.Logger().Fatal(err)
			}
			csvFolder, err := secretDriver.GetCsvFolder(user)
			if err != nil {
				c.Logger().Fatal(err)
			}
			configs[user] = Zaim{
				OAuthConfig: &oauth1.Config{
					ConsumerKey:    consumerKey,
					ConsumerSecret: consumerSecret,
					CallbackURL:    fmt.Sprintf("%s/auth/callback", os.Getenv("HOST")),
					Endpoint: oauth1.Endpoint{
						RequestTokenURL: requestURL,
						AuthorizeURL:    authorizeURL,
						AccessTokenURL:  accessURL,
					},
				},
				OAuthToken: OAuthToken{
					Token:  oauthToken,
					Secret: oauthSecret,
				},
				CsvFolder: csvFolder,
			}
		}
		c.Logger().Debugf("configs: %v", configs)
		req.Body = io.NopCloser(bytes.NewBuffer(b))
		return next(&CustomContext{c, configs})
	}
}
