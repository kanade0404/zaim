package middlewares

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
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
	PubSubMessage pubsubMessage `json:"message"`
	Subscription  string        `json:"subscription"`
}

type pubsubMessage struct {
	Base64EncodedData string `json:"data"`
	MessageID         string `json:"messageId"`
	Message_ID        string `json:"message_id"`
	PublishTime       string `json:"publishTime"`
	Publish_time      string `json:"publish_time"`
}

type requestData struct {
	Users []string `json:"users"`
}

func Context(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().Info("start middleware/context")
		defer c.Logger().Info("end middleware/context")
		ctx := c.Request().Context()
		var body body
		b, err := io.ReadAll(c.Request().Body)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, err)
		}
		c.Logger().Info(string(b))
		if err := json.Unmarshal(b, &body); err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, err)
		}
		b64, err := base64.StdEncoding.DecodeString(body.PubSubMessage.Base64EncodedData)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, err)
		}
		var data requestData
		if err := json.Unmarshal(b64, &data); err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, err)
		}
		users := data.Users
		secretDriver, err := secret_manager.NewDriver(ctx)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, err)
		}
		configs := make(map[string]Zaim, len(users))
		for _, user := range users {
			consumerKey, err := secretDriver.GetConsumerKey(user)
			if err != nil {
				c.Logger().Error(err)
				return c.JSON(http.StatusBadRequest, err)
			}
			consumerSecret, err := secretDriver.GetConsumerSecret(user)
			if err != nil {
				c.Logger().Error(err)
				return c.JSON(http.StatusBadRequest, err)
			}
			oauthToken, err := secretDriver.GetOAuthToken(user)
			if err != nil {
				c.Logger().Error(err)
				return c.JSON(http.StatusBadRequest, err)
			}
			oauthSecret, err := secretDriver.GetOAuthSecret(user)
			if err != nil {
				c.Logger().Error(err)
				return c.JSON(http.StatusBadRequest, err)
			}
			csvFolder, err := secretDriver.GetCsvFolder(user)
			if err != nil {
				c.Logger().Error(err)
				return c.JSON(http.StatusBadRequest, err)
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
		// 後続でstreamを消費しないようにする
		c.Request().Body = io.NopCloser(bytes.NewBuffer(b))
		c.Logger().Debugf("configs: %v", configs)
		return next(&CustomContext{c, configs})
	}
}
