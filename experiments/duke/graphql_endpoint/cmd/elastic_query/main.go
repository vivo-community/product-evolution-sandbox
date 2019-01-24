package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/OIT-ads-web/graphql_endpoint"
	"github.com/davecgh/go-spew/spew"
	"github.com/olivere/elastic"
	"os"
	"time"
)

func listAll(index string) {
	ctx := context.Background()
	client := graphql_endpoint.GetClient()
	q := elastic.NewMatchAllQuery()

	searchResult, err := client.Search().
		Index(index).
		//Type().
		Query(q).
		From(100).
		Size(100).
		Pretty(true).
		// Timeout("1000ms"). or
		// Timeout(1000).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	spew.Println(searchResult.TotalHits())

	/*
	for _, hit := range searchResult.Hits.Hits {
		person := graphql_endpoint.Person{}
		err := json.Unmarshal(*hit.Source, &person)
		if err != nil {
			panic(err)
		}
		spew.Println(person)
	}
	*/
}

var conf graphql_endpoint.Config

func idQuery() {
	q := elastic.NewIdsQuery("person").Ids("per4774112", "per8608642") //.QueryName("my_query")
	ctx := context.Background()
	client := graphql_endpoint.GetClient()

	searchResult, err := client.Search().
		Index("people").
		//Type().
		Query(q).
		From(0).
		Size(1000).
		Pretty(true).
		// Timeout("1000ms"). or
		// Timeout(1000).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	for _, hit := range searchResult.Hits.Hits {
		person := graphql_endpoint.Person{}
		err := json.Unmarshal(*hit.Source, &person)
		if err != nil {
			panic(err)
		}
		spew.Println(person)
	}
}

func findOne(id string) {
	ctx := context.Background()
	client := graphql_endpoint.GetClient()

	get1, err := client.Get().
		Index("people").
		Id(id).
		Do(ctx)

	switch {
	case elastic.IsNotFound(err):
		fmt.Println("404 not found")
	case elastic.IsConnErr(err):
		fmt.Println("connectino error")
	case elastic.IsTimeout(err):
		fmt.Println("timeout")
	case err != nil:
		panic(err)
	}

	var person = graphql_endpoint.Person{}
	err = json.Unmarshal(*get1.Source, &person)
	if err != nil {
		panic(err)
	}

	spew.Println(person)
}

func main() {
	start := time.Now()
	var err error
	var configFile string
	flag.StringVar(&configFile, "config", "./config.toml", "a config filename")

	typeName := flag.String("type", "people", "type of records to query")
	findId := flag.String("id", "per7045252", "id to find")
	flag.Parse()

	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		fmt.Println("could not find config file, use -c option")
		os.Exit(1)
	}

	// NOTE: elastic client is supposed to be long-lived
	// see https://github.com/olivere/elastic/blob/release-branch.v6/client.go
	graphql_endpoint.Client, err = elastic.NewClient(elastic.SetURL(conf.Elastic.Url), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	fmt.Println(*typeName)
	listAll(*typeName)

	fmt.Println("******************")
	findOne(*findId)
	defer graphql_endpoint.Client.Stop()

	fmt.Println("*****************")
	idQuery()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
