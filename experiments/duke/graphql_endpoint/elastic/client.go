package elastic

import (
	"github.com/olivere/elastic"
)

var Client *elastic.Client

func GetClient() *elastic.Client {
	return Client
}

// NOTE: elastic client is supposed to be long-lived
// see https://github.com/olivere/elastic/blob/release-branch.v6/client.go
//client, err = elastic.NewClient(elastic.SetURL(conf.Elastic.Url), elastic.SetSniff(false))
//elastic.Client, err = elastic.NewClient(elastic.SetURL(conf.Elastic.Url), elastic.SetSniff(false))

func MakeClient(url string) (bool, error) {
	// establishing a 'global' client
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
    Client = client

	return true, err
	/*if err != nil {
		panic(err)
	}
	*/
	// establishing a 'global' client
	//Client = client
}

