package elastic

import (
	"net/http"

	"github.com/olivere/elastic"
)

var Client *elastic.Client

func GetClient() *elastic.Client {
	return Client
}

func MakeClient(url string) error {
	// establishing a 'global' client
	client, err := elastic.NewClient(elastic.SetURL(url),
		elastic.SetSniff(false))

	// NOTE: this is establishing a global client because the elastic client is
	// supposed to be long-lived
	// see https://github.com/olivere/elastic/blob/release-branch.v6/client.go
	Client = client
	return err
}

func MakeClientDebug(url string, httpClient *http.Client) error {
	// establishing a 'global' client
	client, err := elastic.NewClient(elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHttpClient(httpClient))

	// NOTE: this is establishing a global client because the elastic client is
	// supposed to be long-lived
	// see https://github.com/olivere/elastic/blob/release-branch.v6/client.go
	Client = client
	return err
}
