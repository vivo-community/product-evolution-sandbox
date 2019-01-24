package graphql_endpoint

import (
	//"context"
	//"github.com/davecgh/go-spew/spew"
	"github.com/olivere/elastic"
)


var Client *elastic.Client

func GetClient() *elastic.Client {
	return Client
}

/*
func listAll(index string) {
	ctx := context.Background()
	client := GetClient()
	q := elastic.NewMatchAllQuery()

	searchResult, err := client.Search().
		Index(index).
		Query(q).
		From(0).
		Size(1000).
		Pretty(true).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	return searchResult.Hits
}

func idQuery() {
	q := elastic.NewIdsQuery("person").Ids("per4774112", "per8608642")//.QueryName("my_query")
	ctx := context.Background()
	client := graphql_endpoint.GetClient()

	searchResult, err := client.Search().
		Index("people").
		//Type().
		Query(q).
		From(0).
		Size(1000).
		Pretty(true).
		Timeout("1000ms").
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	return searchResult.Hits
}
*/
