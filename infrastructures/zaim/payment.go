package zaim

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

type PaymentParameter struct {
	CategoryID    string `json:"category_id"`
	GenreID       string `json:"genre_id"`
	Amount        string `json:"amount"`
	Date          string `json:"date"`
	FromAccountID string `json:"from_account_id"`
	Name          string `json:"name"`
	Place         string `json:"place"`
	Comment       string `json:"comment"`
}

type Money struct {
	ID       int     `json:"id"`
	PlaceUID *string `json:"place_uid"`
	Modified string  `json:"modified"`
}
type User struct {
	InputCount   int    `json:"input_count"`
	RepeatCount  int    `json:"repeat_count"`
	DayCount     int    `json:"day_count"`
	DataModified string `json:"data_modified"`
}
type Place struct {
	ID                int    `json:"id"`
	UserID            int    `json:"user_id"`
	GenreID           int    `json:"genre_id"`
	AccountID         int    `json:"account_id"`
	TransferAccountID int    `json:"transfer_account_id"`
	Mode              string `json:"mode"`
	PlaceUID          string `json:"place_uid"`
	Service           string `json:"service"`
	Name              string `json:"name"`
	OriginalName      string `json:"original_name"`
	Tel               string `json:"tel"`
	Count             int    `json:"count"`
	PlacePatternID    int    `json:"place_pattern_id"`
	CalcFlag          int    `json:"calc_flag"`
	EditFlag          int    `json:"edit_flag"`
	Active            int    `json:"active"`
	Modified          string `json:"modified"`
	Created           string `json:"created"`
}
type PaymentResponse struct {
	Stamps    *string  `json:"stamps"`
	Banners   []string `json:"banners"`
	Money     Money    `json:"money"`
	User      User     `json:"user"`
	Place     Place    `json:"place"`
	Requested int      `json:"requested"`
}

func (z *ZaimClient) CreatePayment(param PaymentParameter) (PaymentResponse, error) {
	obj, err := json.Marshal(param)
	if err != nil {
		return PaymentResponse{}, err
	}
	var paramMap map[string]interface{}
	if err := json.Unmarshal(obj, &paramMap); err != nil {
		return PaymentResponse{}, err
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
			return PaymentResponse{}, fmt.Errorf("unsupported type for key %s", key)
		}
	}
	fmt.Println(values)
	res, err := z.exec("POST", "home/money/payment", values)
	if err != nil {
		return PaymentResponse{}, err
	}
	defer res.Body.Close()
	var r PaymentResponse
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return PaymentResponse{}, err
	}
	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		return PaymentResponse{}, fmt.Errorf("status code is not 2xx. code: %d, body: %v", res.StatusCode, string(b))
	}
	if err := json.Unmarshal(b, &r); err != nil {
		return PaymentResponse{}, err
	}
	return r, nil
}
