package main

import (
	"flag"
	"fmt"
	ge "github.com/OIT-ads-web/graphql_endpoint"
	"github.com/OIT-ads-web/graphql_endpoint/elastic"
	"github.com/OIT-ads-web/graphql_endpoint/examples"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)


func example1() {
	examples.ExampleAggregations()
}

// just a few simple functions to print out data
func main() {
	var conf ge.Config
	start := time.Now()

	viper.SetDefault("elastic.url", "http://localhost:9200")

	if os.Getenv("ENVIRONMENT") == "development" {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")

		value, exists := os.LookupEnv("CONFIG_PATH")
		if exists {
			viper.AddConfigPath(value)
		}

		viper.ReadInConfig()
	} else {
		replacer := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(replacer)
		viper.BindEnv("elastic.url")
	}

	if err := viper.Unmarshal(&conf); err != nil {
		fmt.Printf("could not establish read into conf structure %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("trying to connect to elastic at %s\n", conf.Elastic.Url)

	if err := elastic.MakeClient(conf.Elastic.Url); err != nil {
		fmt.Printf("could not establish elastic client %s\n", err)
		os.Exit(1)
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	fmt.Println("******* aggregations ****")
	// cmd switch for example?
	example1()

	defer elastic.Client.Stop()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
