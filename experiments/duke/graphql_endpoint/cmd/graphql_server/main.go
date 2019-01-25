package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	ge "github.com/OIT-ads-web/graphql_endpoint"
	"github.com/OIT-ads-web/graphql_endpoint/elastic"
	"github.com/OIT-ads-web/graphql_endpoint/graphql"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

var conf ge.Config

func main() {
	var configFile string

	log.SetOutput(os.Stdout)

	flag.StringVar(&configFile, "config", "./config.toml", "a config filename")

	flag.Parse()

	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		fmt.Println("could not find config file, use -c option")
		os.Exit(1)
	}

	if err := elastic.MakeClient(conf.Elastic.Url); err != nil {
		fmt.Printf("could not establish elastic client %s\n", err)
		os.Exit(1)
	}

	c := cors.New(cors.Options{
		AllowCredentials: true,
	})

	handler := graphql.MakeHandler()
	http.Handle("/graphql", c.Handler(handler))

	// NOTE: if not configured this would default to 0
	var port = 9001
	if conf.Graphql.Port > 0 {
		port = conf.Graphql.Port
	}

	portConfig := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portConfig, nil)
}
