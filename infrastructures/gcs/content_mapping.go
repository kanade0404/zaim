package gcs

import (
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
