package graphql_endpoint

import (
	"github.com/olivere/elastic"
)

var Client *elastic.Client

func GetClient() *elastic.Client {
	return Client
}
