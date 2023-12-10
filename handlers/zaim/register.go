package zaim

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"time"
	"zaim/handlers"
	"zaim/infrastructures/zaim"
	"zaim/usecases"
)

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
	RunAt  *string  `json:"run_at"`
	DryRun bool     `json:"dry_run"`
	Users  []string `json:"users"`
}

type RegisterResponse struct {
	Responses [][]zaim.PaymentResponse `json:"responses"`
}

func Register(c echo.Context) error {
	var (
		body  body
		data  requestData
		runAt synchro.Time[tz.AsiaTokyo]
	)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.Logger().Error(err)
		}
	}(c.Request().Body)
	if err := c.Bind(&body); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	b64, err := base64.StdEncoding.DecodeString(body.PubSubMessage.Base64EncodedData)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := json.Unmarshal(b64, &data); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	if data.RunAt == nil {
		runAt = synchro.Now[tz.AsiaTokyo]()
	} else {
		var err error
		runAt, err = synchro.Parse[tz.AsiaTokyo](time.DateOnly, *(data.RunAt))
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, err)
		}
	}
	// jstNowを先月の1日にする
	jstLastMonth := runAt.AddDate(0, -1, -runAt.Day()+1)
	res, err := usecases.RegisterMonthlyTransactions(c, jstLastMonth.StdTime(), data.DryRun)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
			Error: err,
		})
	}
	return c.JSON(http.StatusOK, RegisterResponse{Responses: res})
}
