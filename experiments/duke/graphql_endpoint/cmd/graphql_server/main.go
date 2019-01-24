package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/OIT-ads-web/graphql_endpoint"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/olivere/elastic"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

/*
type Config struct {
	Elastic elasticSearch `toml:"elastic"`
	Graphql graphqlServer `toml:"graphql"`
}

type elasticSearch struct {
	Url string
}

type graphqlServer struct {
	Port int
}
*/

var conf graphql_endpoint.Config

func main() {
	var err error
	var configFile string

	log.SetOutput(os.Stdout)

	flag.StringVar(&configFile, "config", "./config.toml", "a config filename")

	flag.Parse()

	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		fmt.Println("could not find config file, use -c option")
		os.Exit(1)
	}

	// NOTE: elastic client is supposed to be long-lived
	// see https://github.com/olivere/elastic/blob/release-branch.v6/client.go
	//client, err = elastic.NewClient(elastic.SetURL(conf.Elastic.Url), elastic.SetSniff(false))
	graphql_endpoint.Client, err = elastic.NewClient(elastic.SetURL(conf.Elastic.Url), elastic.SetSniff(false))

	if err != nil {
		panic(err)
	}

	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql_endpoint.RootQuery,
	})

	c := cors.New(cors.Options{
		AllowCredentials: true,
	})

	h := handler.New(&handler.Config{
		Schema:   &schema,
		GraphiQL: true,
		Pretty:   true,
	})

	http.Handle("/graphql", c.Handler(h))

	// NOTE: if not configured this would default to 0
	var port = 9001
	if conf.Graphql.Port > 0 {
		port = conf.Graphql.Port
	}

	portConfig := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portConfig, nil)
}
