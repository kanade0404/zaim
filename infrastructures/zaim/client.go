package zaim

import (
	"bytes"
	"fmt"
	"github.com/dghubble/oauth1"
	"net/http"
	"net/url"
)

type ZaimClient struct {
	client *http.Client
}

func NewClient(consumerKey, consumerSecret, token, tokenSecret string) (*ZaimClient, error) {
	oauthConfig := oauth1.NewConfig(consumerKey, consumerSecret)
	oauthToken := oauth1.NewToken(token, tokenSecret)
	httpClient := oauthConfig.Client(oauth1.NoContext, oauthToken)
	return &ZaimClient{httpClient}, nil
}

func (z *ZaimClient) get(path string, values url.Values) (*http.Response, error) {
	return z.exec("GET", path, values)
}

func (z *ZaimClient) post(path string, values url.Values) (*http.Response, error) {
	return z.exec("POST", path, values)
}

func (z *ZaimClient) exec(method, path string, values url.Values) (*http.Response, error) {
	if method == "" {
		method = "GET"
	}
	var (
		req *http.Request
		err error
	)
	if method == "GET" {
		req, err = http.NewRequest(method, fmt.Sprintf("https://api.zaim.net/v2/%s?%s", path, values.Encode()), nil)
	} else {
		req, err = http.NewRequest(method, fmt.Sprintf("https://api.zaim.net/v2/%s", path), bytes.NewBufferString(values.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	if err != nil {
		return nil, err
	}
	return z.client.Do(req)
}
