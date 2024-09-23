package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dghubble/oauth1"
	"net/http"
	"net/url"
)

type ZaimDriver struct {
	client *http.Client
}

const baseURL = "https://api.zaim.net/v2"

func NewZaimDriver(consumerKey, consumerSecret, token, tokenSecret string) (ZaimDriver, error) {
	oauthConfig := oauth1.NewConfig(consumerKey, consumerSecret)
	oauthToken := oauth1.NewToken(token, tokenSecret)
	httpClient := oauthConfig.Client(oauth1.NoContext, oauthToken)
	return ZaimDriver{httpClient}, nil
}

func (z ZaimDriver) Get(path string, values interface{}) (*http.Response, error) {
	return z.exec("GET", path, values)
}

func (z ZaimDriver) Post(path string, values interface{}) (*http.Response, error) {
	return z.exec("POST", path, values)
}

func (z ZaimDriver) Delete(path string, values interface{}) (*http.Response, error) {
	return z.exec("DELETE", path, values)
}

func (z ZaimDriver) exec(method, path string, params interface{}) (*http.Response, error) {
	if method == "" {
		method = "GET"
	}
	var paramMap map[string]interface{}
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
	req, err := createRequest(method, path, values)
	if err != nil {
		return nil, err
	}
	return z.client.Do(req)
}

func createRequest(method, path string, values url.Values) (*http.Request, error) {
	switch method {
	case "GET":
		return http.NewRequest(method, fmt.Sprintf("%s/%s?%s", baseURL, path, values.Encode()), nil)
	case "DELETE":
		return http.NewRequest(method, fmt.Sprintf("%s/%s?%s", baseURL, path, values.Encode()), nil)
	case "POST", "PUT":
		req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", baseURL, path), bytes.NewBufferString(values.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		return req, nil
	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}
}
