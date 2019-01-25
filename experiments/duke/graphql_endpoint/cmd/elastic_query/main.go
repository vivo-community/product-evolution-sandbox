package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	ge "github.com/OIT-ads-web/graphql_endpoint"
	"github.com/OIT-ads-web/graphql_endpoint/elastic"
	"os"
	"time"
)

func listAll(index string) {
	elastic.ListAll(index)
}

func idQuery() {
	elastic.IdQuery("person", []string{"per4774112", "per8608642"})
}

func findOne(id string) {
	elastic.FindOne(id)
}

var conf ge.Config

func main() {
	start := time.Now()
	var configFile string
	flag.StringVar(&configFile, "config", "./config.toml", "a config filename")

	typeName := flag.String("type", "people", "type of records to query")
	findId := flag.String("id", "per7045252", "id to find")
	flag.Parse()

	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		fmt.Println("could not find config file, use -c option")
		os.Exit(1)
	}

	if _, err := elastic.MakeClient(conf.Elastic.Url); err != nil {
		fmt.Printf("could not establish elastic client %s\n", err)
		os.Exit(1)
	}

	fmt.Println(*typeName)
	listAll(*typeName)

	fmt.Println("******************")
	findOne(*findId)
	defer elastic.Client.Stop()

	fmt.Println("*****************")
	idQuery()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
