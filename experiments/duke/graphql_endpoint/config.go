package graphql_endpoint

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

