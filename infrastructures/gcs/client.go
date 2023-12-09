package gcs

import (
	"bufio"
	"cloud.google.com/go/storage"
	"context"
	"github.com/labstack/gommon/log"
	"io"
)

type Driver struct {
	Client *storage.Client
	ctx    context.Context
}

func NewDriver(ctx context.Context) (*Driver, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Driver{client, ctx}, nil
}
func (d *Driver) Upload(bucketName, objectPath string, obj io.Reader) error {
	writer := d.Client.Bucket(bucketName).Object(objectPath).NewWriter(d.ctx)
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

func (d *Driver) Download(bucketName, fileName string) ([]byte, error) {
	bucket := d.Client.Bucket(bucketName)
	rc, err := bucket.Object(fileName).NewReader(d.ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
}
