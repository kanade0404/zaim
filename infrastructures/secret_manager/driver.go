package secret_manager

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"fmt"
	"os"
)

type Driver struct {
	ctx    context.Context
	client *secretmanager.Client
}

func NewDriver(ctx context.Context) (*Driver, error) {
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Driver{
		ctx:    ctx,
		client: c,
	}, nil
}

func (d *Driver) accessSecretVersionRequest(name string) (string, error) {
	res, err := d.client.AccessSecretVersion(d.ctx,
		&secretmanagerpb.AccessSecretVersionRequest{
			Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", os.Getenv("PROJECT_ID"), name),
		})
	if err != nil {
		return "", err
	}
	return string(res.GetPayload().Data), nil
}
