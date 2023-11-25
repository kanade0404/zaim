package zaim

import (
	"encoding/json"
	"io"
)

type Account struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Modified        string `json:"modified"`
	Sort            int    `json:"sort"`
	Active          int    `json:"active"`
	WebsiteID       int    `json:"website_id"`
	ParentAccountID int    `json:"parent_account_id"`
}

func (z *ZaimClient) ListActiveAccount() ([]Account, error) {
	res, err := z.get("home/account", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var r struct {
		Accounts  []Account `json:"accounts"`
		Requested int       `json:"requested"`
	}
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	var results []Account
	for i := range r.Accounts {
		if r.Accounts[i].Active == 1 {
			results = append(results, r.Accounts[i])
		}
	}
	return results, nil
}
