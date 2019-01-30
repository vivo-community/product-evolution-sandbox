package main

import (
	"fmt"
	ge "github.com/OIT-ads-web/graphql_endpoint"
	"github.com/OIT-ads-web/graphql_endpoint/elastic"
	"github.com/OIT-ads-web/graphql_endpoint/graphql"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"strings"
)

var conf ge.Config

func main() {
	log.SetOutput(os.Stdout)

	viper.SetDefault("elastic.url", "http://localhost:9200")
	viper.SetDefault("graphql.port", "9001")

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

	if err := viper.Unmarshal(&conf); err != nil {
		fmt.Printf("could not establish read into conf structure %s\n", err)
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

	port := conf.Graphql.Port
	portConfig := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portConfig, nil)
}
