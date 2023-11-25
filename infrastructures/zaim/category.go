package zaim

import (
	"encoding/json"
	"io"
)

type Category struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Mode             string `json:"mode"`
	Sort             int    `json:"sort"`
	ParentCategoryID int    `json:"parent_category_id"`
	Active           int    `json:"active"`
	Modified         string `json:"modified"`
}

func (z *ZaimClient) ListActiveCategory() ([]Category, error) {
	res, err := z.get("home/category", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var r struct {
		Categories []Category `json:"categories"`
		Requested  int        `json:"requested"`
	}
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	var results []Category
	for i := range r.Categories {
		if r.Categories[i].Active == 1 {
			results = append(results, r.Categories[i])
		}

	}
	return results, nil
}
