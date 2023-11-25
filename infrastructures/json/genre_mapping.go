package json

import (
	"encoding/json"
	"fmt"
	"os"
)

type GenreMappings map[string]GenreContentMappings
type GenreContentMappings map[string]GenreCategory
type GenreCategory struct {
	CategoryID int `json:"category_id"`
	GenreID    int `json:"genre_id"`
}

func GetGenreMappingByUserName(userName string) (GenreContentMappings, error) {
	f, err := os.Open("./database/genre_mapping.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var genreMappings GenreMappings
	if err := json.NewDecoder(f).Decode(&genreMappings); err != nil {
		return nil, err
	}
	if genreMapping, ok := genreMappings[userName]; ok {
		return genreMapping, nil
	} else {
		return nil, fmt.Errorf("genre mapping not found: %s", userName)
	}
}
