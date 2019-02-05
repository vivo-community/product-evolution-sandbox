package graphql

import (
	"github.com/graphql-go/handler"
)

func MakeHandler() *handler.Handler {
	schema := MakeSchema()
	h := handler.New(&handler.Config{
		Schema:   &schema,
		GraphiQL: true,
		Pretty:   true,
	})
	return h
}
