package gcs

import (
	"bytes"
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

func PutGenreMapping(ctx context.Context, genreMappings GenreMappings) error {
	gcsDriver, err := NewDriver(ctx)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(genreMappings); err != nil {
		return err
	}
	if err := gcsDriver.Upload("zaim", "genre_mappings.json", buf); err != nil {
		return err
	}
	return nil
}
