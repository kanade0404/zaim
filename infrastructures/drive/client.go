package drive

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type DriveClient struct {
	Client *http.Client
	ctx    context.Context
}

func NewClient(ctx context.Context) (*DriveClient, error) {
	p, err := filepath.Abs("../../credentials/gdrive_credential.json")
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		return nil, err
	}
	tokenFile, err := filepath.Abs("../../credentials/token.json")
	if err != nil {
		return nil, err
	}
	token, err := tokenFromFile(tokenFile)
	if err != nil {
		token, err = getTokenFromWeb(ctx, config)
		if err != nil {
			return nil, err
		}
		if err := saveToken(tokenFile, token); err != nil {
			return nil, err
		}
	}
	return &DriveClient{config.Client(ctx, token), ctx}, nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web %v", err)
	}
	return tok, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}(f)
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %v", err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}(f)
	return json.NewEncoder(f).Encode(token)
}

type CreateChannelPayload struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Address string `json:"address"`
}

func (d *DriveClient) CreateStartPageToken() (*drive.StartPageToken, error) {
	srv, err := drive.NewService(d.ctx, option.WithHTTPClient(d.Client))
	if err != nil {
		return nil, err
	}
	token, err := srv.Changes.GetStartPageToken().Do()
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (d *DriveClient) CreateWatch(folderID string, channel *drive.Channel) (*drive.Channel, error) {
	srv, err := drive.NewService(d.ctx, option.WithHTTPClient(d.Client))
	if err != nil {
		return nil, err
	}
	resp, err := srv.Files.Watch(folderID, channel).Do()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
