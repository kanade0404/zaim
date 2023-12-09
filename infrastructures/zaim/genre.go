package zaim

import (
	"encoding/json"
	"io"
)

type Genre struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Sort          int    `json:"sort"`
	Active        int    `json:"active"`
	CategoryID    int    `json:"category_id"`
	ParentGenreID int    `json:"parent_genre_id"`
	Modified      string `json:"modified"`
}

func (z *Client) ListActiveGenre() ([]Genre, error) {
	res, err := z.get("home/genre", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var r struct {
		Genres    []Genre `json:"genres"`
		Requested int     `json:"requested"`
	}
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	var results []Genre
	for i := range r.Genres {
		if r.Genres[i].Active == 1 {
			results = append(results, r.Genres[i])
		}
	}
	return results, nil
}
