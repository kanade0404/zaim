package json

import (
	"encoding/json"
	"fmt"
	"os"
)

type Category struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Mode             string `json:"mode"`
	Sort             int    `json:"sort"`
	ParentCategoryID *int   `json:"parent_category_id"`
	Active           int    `json:"active"`
	Modified         string `json:"modified"`
}

type Categories map[string][]Category

func GetCategoryByUserName(userName string) ([]Category, error) {
	f, err := os.Open("./database/category.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var categories Categories
	if err := json.NewDecoder(f).Decode(&categories); err != nil {
		return nil, err
	}
	if category, ok := categories[userName]; ok {
		return category, nil
	} else {
		return nil, fmt.Errorf("category not found: %s", userName)
	}

}
