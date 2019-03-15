package graphql

import (
	"github.com/graphql-go/graphql"
)

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"personList":      GetPeople,
		"person":          GetPerson,
		"publicationList": GetPublications,
		"grantList":       GetGrants,
	},
})

var GetPerson = &graphql.Field{
	Type:        person,
	Description: "Get Person",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: personResolver,
}

var filterObject *graphql.InputObject = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "Filter",
	Fields: graphql.InputObjectConfigFieldMap{
		"limit": &graphql.InputObjectFieldConfig{
			Type:         graphql.Int,
			DefaultValue: 100,
		},
		"offset": &graphql.InputObjectFieldConfig{
			Type:         graphql.Int,
			DefaultValue: 0,
		},
	},
})

var GetPeople = &graphql.Field{
	Type:        personList,
	Description: "Get all people",
	Args: graphql.FieldConfigArgument{
		"filter": &graphql.ArgumentConfig{Type: filterObject},
	},
	Resolve: peopleResolver,
}

var GetPublications = &graphql.Field{
	Type:        publicationList,
	Description: "Get all publications",
	Args: graphql.FieldConfigArgument{
		"size": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 100},
		"from": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
	},
	Resolve: publicationResolver,
}

var GetGrants = &graphql.Field{
	Type:        grantList,
	Description: "Get all grants",
	Args: graphql.FieldConfigArgument{
		"size": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 100},
		"from": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
	},
	Resolve: grantResolver,
}
