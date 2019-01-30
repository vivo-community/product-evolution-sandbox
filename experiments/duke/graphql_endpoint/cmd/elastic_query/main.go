package main

import (
	"flag"
	"fmt"
	ge "github.com/OIT-ads-web/graphql_endpoint"
	"github.com/OIT-ads-web/graphql_endpoint/elastic"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"time"
	"strings"
)

func listAll(index string) {
	elastic.ListAll(index)
}

func idQuery() {
	elastic.IdQuery("people", []string{"per4774112", "per8608642"})
}

func findOne(index string, id string) {
	elastic.FindOne(index, id)
}

var conf ge.Config

// just a few simple functions to print out data
func main() {
	start := time.Now()

	viper.SetDefault("elastic.url", "http://localhost:9200")

	if os.Getenv("ENVIRONMENT") == "development" {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")
		viper.ReadInConfig()
	} else {
		replacer := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(replacer)
		viper.AutomaticEnv()
	}

	fmt.Printf("trying to connect to elastic at %s\n", viper.GetString("elastic.url"))

	if err := elastic.MakeClient(viper.GetString("elastic.url")); err != nil {
		fmt.Printf("could not establish elastic client %s\n", err)
		os.Exit(1)
	}

	flag.String("type", "people", "type of records to query")
	flag.String("id", "per7045252", "id to find")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	fmt.Println(viper.GetString("type"))
	listAll(viper.GetString("type"))

	fmt.Println("******************")
	findOne(viper.GetString("type"), viper.GetString("id"))

	defer elastic.Client.Stop()

	fmt.Println("*****************")

	idQuery()
	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
