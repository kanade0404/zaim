package secret_manager

import (
	"fmt"
)

func (d *Driver) GetConsumerKey(userName string) (string, error) {
	return d.accessSecretVersionRequest(fmt.Sprintf("zaim-consumer-key-%s", userName))
}

func (d *Driver) GetConsumerSecret(userName string) (string, error) {
	return d.accessSecretVersionRequest(fmt.Sprintf("zaim-consumer-secret-%s", userName))
}

func (d *Driver) GetCSVFolder(userName string) (string, error) {
	return d.accessSecretVersionRequest(fmt.Sprintf("csv-folder-%s", userName))
}

func (d *Driver) GetOAuthToken(userName string) (string, error) {
	return d.accessSecretVersionRequest(fmt.Sprintf("zaim-oauth-token-%s", userName))
}

func (d *Driver) GetOAuthSecret(userName string) (string, error) {
	return d.accessSecretVersionRequest(fmt.Sprintf("zaim-oauth-secret-%s", userName))
}
func (d *Driver) GetCsvFolder(userName string) (string, error) {
	return d.accessSecretVersionRequest(fmt.Sprintf("csv-folder-%s", userName))
}
