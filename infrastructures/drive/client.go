package drive

import (
	"context"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"net/http"
)

type DriveClient struct {
	Client *http.Client
	ctx    context.Context
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
