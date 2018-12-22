package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
	"github.com/BurntSushi/toml"
	"github.com/OIT-ads-web/graphql_endpoint"
	"github.com/OIT-ads-web/graphql_endpoint/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/olivere/elastic"
)

var client *elastic.Client

func GetClient() *elastic.Client {
	return client
}

func listAll(index string) {
	ctx := context.Background()
	client = GetClient()
	q := elastic.NewMatchAllQuery()

	searchResult, err := client.Search().
		Index(index).
		//Type().
		Query(q).
		From(0).
		Size(1000).
		//Pretty(true).
		// Timeout("1000ms"). or
		// Timeout(1000).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	//TotalHits()

	for _, hit := range searchResult.Hits.Hits {
		person := models.Person{}
		err := json.Unmarshal(*hit.Source, &person)
		if err != nil {
			panic(err)
		}
		spew.Println(person)
	}
}

var conf graphql_endpoint.Config

func idQuery() {
	q := elastic.NewIdsQuery("person").Ids("per4774112", "per8608642")//.QueryName("my_query")
	/*
	src, err := q.Source()
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		log.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	fmt.Println(got)
    */
	ctx := context.Background()
	client = GetClient()

	searchResult, err := client.Search().
		Index("people").
		//Type().
		Query(q).
		From(0).
		Size(1000).
		//Pretty(true).
		// Timeout("1000ms"). or
		// Timeout(1000).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	for _, hit := range searchResult.Hits.Hits {
		person := models.Person{}
		err := json.Unmarshal(*hit.Source, &person)
		if err != nil {
			panic(err)
		}
		spew.Println(person)
	}
}

func findOne(id string) {
	ctx := context.Background()
	client = GetClient()

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

	var person = models.Person{}
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
	client, err = elastic.NewClient(elastic.SetURL(conf.Elastic.Url))
	if err != nil {
		panic(err)
	}

	fmt.Println(*typeName)
	//listAll(*typeName)

	fmt.Println("******************")
	findOne(*findId)
	defer client.Stop()

	fmt.Println("*****************")
	idQuery()
	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
