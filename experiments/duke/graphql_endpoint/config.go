package graphql_endpoint

type Config struct {
	Elastic  elasticSearch `toml:"elastic"`
}

type elasticSearch struct {
	Url string
}


