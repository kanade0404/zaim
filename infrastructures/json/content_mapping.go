package json

import (
	"encoding/json"
	"fmt"
	"os"
)

type ContentMappings map[string]ContentContentMappings
type ContentContentMappings map[string]GenreCategory
type ContentCategory struct {
	CategoryID int `json:"category_id"`
	GenreID    int `json:"genre_id"`
}

func GetContentMappingByUserName(userName string) (ContentContentMappings, error) {
	f, err := os.Open("./database/content_mapping.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var contentMappings ContentMappings
	if err := json.NewDecoder(f).Decode(&contentMappings); err != nil {
		return nil, err
	}
	if contentMapping, ok := contentMappings[userName]; ok {
		return contentMapping, nil
	} else {
		return nil, fmt.Errorf("content mapping not found: %s", userName)
	}
}
