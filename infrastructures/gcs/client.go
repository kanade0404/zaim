package gcs

import (
	"bufio"
	"cloud.google.com/go/storage"
	"context"
	"github.com/labstack/gommon/log"
	"io"
)

type GcsClient struct {
	Client *storage.Client
	ctx    context.Context
}

func NewClient(ctx context.Context) (*GcsClient, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GcsClient{client, ctx}, nil
}
func (c *GcsClient) Upload(bucketName, objectPath string, obj io.Reader) error {
	writer := c.Client.Bucket(bucketName).Object(objectPath).NewWriter(c.ctx)
	defer func(reader *storage.Writer) {
		err := writer.Close()
		if err != nil {
			log.Error(err)
		}
	}(writer)
	if _, err := io.Copy(writer, bufio.NewReader(obj)); err != nil {
		return err
	}
	return nil
}
