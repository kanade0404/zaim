package zaim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dghubble/oauth1"
	"net/http"
	"net/url"
)

type Client struct {
	client *http.Client
}

func NewClient(consumerKey, consumerSecret, token, tokenSecret string) (*Client, error) {
	oauthConfig := oauth1.NewConfig(consumerKey, consumerSecret)
	oauthToken := oauth1.NewToken(token, tokenSecret)
	httpClient := oauthConfig.Client(oauth1.NoContext, oauthToken)
	return &Client{httpClient}, nil
}

func (z *Client) get(path string, values interface{}) (*http.Response, error) {
	return z.exec("GET", path, values)
}

func (z *Client) post(path string, values interface{}) (*http.Response, error) {
	return z.exec("POST", path, values)
}

func (z *Client) delete(path string, values interface{}) (*http.Response, error) {
	return z.exec("DELETE", path, values)
}

func (z *Client) exec(method, path string, params interface{}) (*http.Response, error) {
	if method == "" {
		method = "GET"
	}
	var (
		paramMap map[string]interface{}
		req      *http.Request
		err      error
	)
	obj, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(obj, &paramMap); err != nil {
		return nil, err
	}
	values := url.Values{}
	for key, value := range paramMap {
		switch v := value.(type) {
		case string:
			values.Add(key, v)
		case float64:
			// JSONは数値をfloat64として扱う
			values.Add(key, fmt.Sprintf("%v", v))
		// 他の型についても必要に応じてケースを追加
		default:
			return nil, fmt.Errorf("unsupported type for key %s", key)
		}
	}
	if method == "GET" {
		req, err = http.NewRequest(method, fmt.Sprintf("https://api.zaim.net/v2/%s?%s", path, values.Encode()), nil)
	} else if method == "DELETE" {
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
