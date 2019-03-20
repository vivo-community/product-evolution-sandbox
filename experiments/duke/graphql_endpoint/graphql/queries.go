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

var PersonFilter *graphql.InputObject = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "PersonFilter",
	Description: "Filter on People List",
	Fields: graphql.InputObjectConfigFieldMap{
		"limit": &graphql.InputObjectFieldConfig{
			Type:         graphql.Int,
			DefaultValue: 100,
		},
		"offset": &graphql.InputObjectFieldConfig{
			Type:         graphql.Int,
			DefaultValue: 0,
		},
		"query": &graphql.InputObjectFieldConfig{
			Type:         graphql.String,
			DefaultValue: "",
		},
	},
})

var GetPeople = &graphql.Field{
	Type:        personList,
	Description: "Get all people",
	Args: graphql.FieldConfigArgument{
		"filter": &graphql.ArgumentConfig{Type: PersonFilter},
	},
	Resolve: peopleResolver,
}

// TODO: very likely a way to avoid the code duplication
var PublicationFilter *graphql.InputObject = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "PublicationFilter",
	Description: "Filter on Publication List",
	Fields: graphql.InputObjectConfigFieldMap{
		"limit": &graphql.InputObjectFieldConfig{
			Type:         graphql.Int,
			DefaultValue: 100,
		},
		"offset": &graphql.InputObjectFieldConfig{
			Type:         graphql.Int,
			DefaultValue: 0,
		},
		"query": &graphql.InputObjectFieldConfig{
			Type:         graphql.String,
			DefaultValue: "",
		},
	},
})

var GetPublications = &graphql.Field{
	Type:        publicationList,
	Description: "Get all publications",
	Args: graphql.FieldConfigArgument{
		"filter": &graphql.ArgumentConfig{Type: PublicationFilter},
	},
	Resolve: publicationResolver,
}

var GrantFilter *graphql.InputObject = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "GrantFilter",
	Description: "Filter on Grant List",
	Fields: graphql.InputObjectConfigFieldMap{
		"limit": &graphql.InputObjectFieldConfig{
			Type:         graphql.Int,
			DefaultValue: 100,
		},
		"offset": &graphql.InputObjectFieldConfig{
			Type:         graphql.Int,
			DefaultValue: 0,
		},
		"query": &graphql.InputObjectFieldConfig{
			Type:         graphql.String,
			DefaultValue: "",
		},
	},
})

var GetGrants = &graphql.Field{
	Type:        grantList,
	Description: "Get all grants",
	Args: graphql.FieldConfigArgument{
		"filter": &graphql.ArgumentConfig{Type: GrantFilter},
	},
	Resolve: grantResolver,
}
