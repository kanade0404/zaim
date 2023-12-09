package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type ContentMappings map[string]ContentContentMappings
type ContentContentMappings map[string]GenreCategory
type ContentCategory struct {
	CategoryID int `json:"category_id"`
	GenreID    int `json:"genre_id"`
}

func GetContentMappingByUserName(ctx context.Context, userName string) (ContentContentMappings, error) {
	gcsDriver, err := NewDriver(ctx)
	if err != nil {
		return nil, err
	}
	data, err := gcsDriver.Download("zaim", "content_mappings.json")
	if err != nil {
		return nil, err
	}
	var contentMappings ContentMappings
	if err := json.Unmarshal(data, &contentMappings); err != nil {
		return nil, err
	}
	if contentContentMappings, ok := contentMappings[userName]; ok {
		return contentContentMappings, nil
	} else {
		return nil, fmt.Errorf("content_mappings not found: %s", userName)
	}
}

func PutContentMappings(ctx context.Context, contentMappings ContentMappings) error {
	gcsDriver, err := NewDriver(ctx)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(contentMappings); err != nil {
		return err
	}
	if err := gcsDriver.Upload("zaim", "content_mappings.json", buf); err != nil {
		return err
	}
	return nil

}
