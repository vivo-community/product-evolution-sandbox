# graphql server in golang

## config 

can either set environmental variables, 

> `export ELASTIC_URL="http://localhost:9200"`
> `export GRAPHQL_PORT="9001"`

or if `set ENVIRONMENT=development` looks for config.toml file
in current directory (see config.toml.example)

## server 

* endpoint on `GRAPHQL_PORT`
* see localhost:<GRAPHQL_PORT>/graphql

