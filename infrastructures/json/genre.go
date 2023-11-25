package json

import (
	"encoding/json"
	"os"
)

type Genres map[string][]Genre
type Genre struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Sort          int    `json:"sort"`
	Active        int    `json:"active"`
	CategoryID    int    `json:"category_id"`
	ParentGenreID int    `json:"parent_genre_id"`
	Modified      string `json:"modified"`
}

func GetGenreByUserName(userName string) ([]Genre, error) {
	f, err := os.Open("./database/genre.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var genres Genres
	if err := json.NewDecoder(f).Decode(&genres); err != nil {
		return nil, err
	}
	if genre, ok := genres[userName]; ok {
		return genre, nil
	} else {
		return nil, err
	}
}
