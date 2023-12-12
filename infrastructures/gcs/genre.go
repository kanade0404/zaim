package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"zaim/infrastructures/zaim"
)

type Genres map[string][]zaim.Genre

func GetGenreByUserName(ctx context.Context, userName string) ([]zaim.Genre, error) {
	gcsDriver, err := NewDriver(ctx)
	if err != nil {
		return nil, err
	}
	data, err := gcsDriver.Download(os.Getenv("BUCKET_NAME"), "genre.json")
	if err != nil {
		return nil, err
	}
	var genres Genres
	if err := json.Unmarshal(data, &genres); err != nil {
		return nil, err
	}
	if genre, ok := genres[userName]; ok {
		return genre, nil
	} else {
		return nil, fmt.Errorf("genre not found: %s", userName)
	}
}

func PutGenre(ctx context.Context, genres Genres) error {
	gcsDriver, err := NewDriver(ctx)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(genres); err != nil {
		return err
	}
	if err := gcsDriver.Upload(os.Getenv("BUCKET_NAME"), "genre.json", buf); err != nil {
		return err
	}
	return nil
}
