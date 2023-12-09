package gcs

import (
	"context"
	"encoding/json"
	"fmt"
)

type GenreMappings map[string]GenreContentMappings
type GenreContentMappings map[string]GenreCategory
type GenreCategory struct {
	CategoryID int `json:"category_id"`
	GenreID    int `json:"genre_id"`
}

func GetGenreMappingByUserName(ctx context.Context, userName string) (GenreContentMappings, error) {
	gcsDriver, err := NewDriver(ctx)
	if err != nil {
		return nil, err
	}
	data, err := gcsDriver.Download("zaim", "genre_mappings.json")
	if err != nil {
		return nil, err
	}
	var genreMappings GenreMappings
	if err := json.Unmarshal(data, &genreMappings); err != nil {
		return nil, err
	}
	if genreContentMappings, ok := genreMappings[userName]; ok {
		return genreContentMappings, nil
	} else {
		return nil, fmt.Errorf("genre_mappings not found: %s", userName)
	}
}
