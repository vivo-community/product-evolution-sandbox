package graphql

//https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
import (
	"github.com/graphql-go/graphql"
)

func MakeSchema() graphql.Schema {
	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: RootQuery,
	})
    return schema
}


