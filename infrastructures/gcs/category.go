package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"zaim/infrastructures/zaim"
)

type Categories map[string][]zaim.Category

func GetCategoryByUserName(ctx context.Context, userName string) ([]zaim.Category, error) {
	gcsDriver, err := NewDriver(ctx)
	if err != nil {
		return nil, err
	}
	data, err := gcsDriver.Download(os.Getenv("BUCKET_NAME"), "category.json")
	if err != nil {
		return nil, err
	}
	var categories Categories
	if err := json.Unmarshal(data, &categories); err != nil {
		return nil, err
	}
	if category, ok := categories[userName]; ok {
		return category, nil
	} else {
		return nil, fmt.Errorf("category not found: %s", userName)
	}
}

func PutCategory(ctx context.Context, categories Categories) error {
	gcsDriver, err := NewDriver(ctx)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(categories); err != nil {
		return err
	}
	if err := gcsDriver.Upload(os.Getenv("BUCKET_NAME"), "category.json", buf); err != nil {
		return err
	}
	return nil
}
