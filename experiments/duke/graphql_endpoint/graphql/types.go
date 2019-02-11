package graphql

import (
	"github.com/graphql-go/graphql"
)

var pageInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PageInfo",
	Fields: graphql.Fields{
		"perPage":    &graphql.Field{Type: graphql.Int},
		"page":       &graphql.Field{Type: graphql.Int},
		"totalPages": &graphql.Field{Type: graphql.Int},
		"count":      &graphql.Field{Type: graphql.Int},
	},
})

var grantType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Grant",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"label":     &graphql.Field{Type: graphql.String},
		"roleName":  &graphql.Field{Type: graphql.String},
		"startDate": &graphql.Field{Type: dateResolutionType},
		"endDate":   &graphql.Field{Type: dateResolutionType},
	},
})

var organizationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Organization",
	Fields: graphql.Fields{
		"id":    &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
	},
})

var educationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Education",
	Fields: graphql.Fields{
		"label": &graphql.Field{Type: graphql.String},
		"org":   &graphql.Field{Type: organizationType},
	},
})

var affiliationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Affiliation",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"label":     &graphql.Field{Type: graphql.String},
		"startDate": &graphql.Field{Type: dateResolutionType},
	},
})

var keywordType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Keyword",
	Fields: graphql.Fields{
		"uri":   &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
	},
})

var extensionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Extension",
	Fields: graphql.Fields{
		"key":   &graphql.Field{Type: graphql.String},
		"value": &graphql.Field{Type: graphql.String},
	},
})

var dateResolutionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DateResolution",
	Fields: graphql.Fields{
		"dateTime":   &graphql.Field{Type: graphql.String},
		"resolution": &graphql.Field{Type: graphql.String},
	},
})

var personNameType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonName",
	Fields: graphql.Fields{
		"firstName":  &graphql.Field{Type: graphql.String},
		"lastName":   &graphql.Field{Type: graphql.String},
		"middleName": &graphql.Field{Type: graphql.String},
	},
})

var personImageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonImage",
	Fields: graphql.Fields{
		"main":      &graphql.Field{Type: graphql.String},
		"thumbnail": &graphql.Field{Type: graphql.String},
	},
})

var typeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Type",
	Fields: graphql.Fields{
		"code":  &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
	},
})

var personIdentifierType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonIdentifier",
	Fields: graphql.Fields{
		"orchid": &graphql.Field{Type: graphql.String},
		"isni":   &graphql.Field{Type: graphql.String},
	},
})

var serviceRoleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ServiceRole",
	Fields: graphql.Fields{
		"uri":          &graphql.Field{Type: graphql.String},
		"label":        &graphql.Field{Type: graphql.String},
		"description":  &graphql.Field{Type: graphql.String},
		"startDate":    &graphql.Field{Type: dateResolutionType},
		"endDate":      &graphql.Field{Type: dateResolutionType},
		"organization": &graphql.Field{Type: organizationType},
		"type":         &graphql.Field{Type: typeType},
	},
})

var contactType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Contact",
	Fields: graphql.Fields{
		"uri":      &graphql.Field{Type: graphql.String},
		"email":    &graphql.Field{Type: graphql.String},
		"phone":    &graphql.Field{Type: graphql.String},
		"location": &graphql.Field{Type: graphql.String},
		"website":  &graphql.Field{Type: graphql.String},
	},
})

var courseTaughtType = graphql.NewObject(graphql.ObjectConfig{
	Name:   "CourseTaught",
	Fields: graphql.Fields{
		"uri":      &graphql.Field{Type: graphql.String},
		"subject":      &graphql.Field{Type: graphql.String},
		"role":      &graphql.Field{Type: graphql.String},
		"courseName":      &graphql.Field{Type: graphql.String},
		"courseNumber":      &graphql.Field{Type: graphql.String},
		"startDate":    &graphql.Field{Type: dateResolutionType},
		"endDate":      &graphql.Field{Type: dateResolutionType},
		"organization": &graphql.Field{Type: organizationType},
		"type":         &graphql.Field{Type: typeType},
	},
})

var publicationIdentifierType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PublicationIdentifier",
	Fields: graphql.Fields{
		"isbn10": &graphql.Field{Type: graphql.String},
		"isbn13": &graphql.Field{Type: graphql.String},
		"pmid":   &graphql.Field{Type: graphql.String},
		"doi":    &graphql.Field{Type: graphql.String},
		"pmcid":  &graphql.Field{Type: graphql.String},
	},
})

var publicationVenueType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PublicationVenue",
	Fields: graphql.Fields{
		"uri":   &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
	},
})

var publicationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Publication",
	Fields: graphql.Fields{
		"id":         &graphql.Field{Type: graphql.String},
		"uri":        &graphql.Field{Type: graphql.String},
		"label":      &graphql.Field{Type: graphql.String},
		"identifier": &graphql.Field{Type: publicationIdentifierType},
		"type":       &graphql.Field{Type: typeType},
		"authorList": &graphql.Field{Type: graphql.String},
		"venue":      &graphql.Field{Type: publicationVenueType},
	},
})

var overviewType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Overview",
	Fields: graphql.Fields{
		"code":  &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
	},
})

var authorshipType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Authorship",
	Fields: graphql.Fields{
		"id":            &graphql.Field{Type: graphql.String},
		"uri":           &graphql.Field{Type: graphql.String},
		"publicationId": &graphql.Field{Type: graphql.String},
		"personId":      &graphql.Field{Type: graphql.String},
		"label":         &graphql.Field{Type: graphql.String},
	},
})

var personType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Person",
	Fields: graphql.Fields{
		"uri":              &graphql.Field{Type: graphql.String},
		"id":               &graphql.Field{Type: graphql.String},
		"sourceId":         &graphql.Field{Type: graphql.String},
		"primaryTitle":     &graphql.Field{Type: graphql.String},
		"name":             &graphql.Field{Type: personNameType},
		"image":            &graphql.Field{Type: personImageType},
		"identifier":       &graphql.Field{Type: personIdentifierType},
		"type":             &graphql.Field{Type: typeType},
		"overviewList":     &graphql.Field{Type: graphql.NewList(overviewType)},
		"keywordList":      &graphql.Field{Type: graphql.NewList(keywordType)},
		"extensionList":    &graphql.Field{Type: graphql.NewList(extensionType)},
		"affliationList":   &graphql.Field{Type: graphql.NewList(affiliationType)},
		"educationList":    &graphql.Field{Type: graphql.NewList(educationType)},
		"serviceRoleList":  &graphql.Field{Type: graphql.NewList(serviceRoleType)},
		//"contactList":      &graphql.Field{Type: graphql.NewList(contactType)},
		"courseTaughtList": &graphql.Field{Type: graphql.NewList(courseTaughtType)},
		// these can be paged, since they involve further queries
		"publicationList": &graphql.Field{
			Type: publicationListType,
			Args: graphql.FieldConfigArgument{
				"size": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 100},
				"from": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: personPublicationResolver,
		},
		"grantList": &graphql.Field{
			Type: grantListType,
			Args: graphql.FieldConfigArgument{
				"size": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 100},
				"from": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: personGrantResolver,
		},
	},
})

var personListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonList",
	Fields: graphql.Fields{
		"results":  &graphql.Field{Type: graphql.NewList(personType)},
		"pageInfo": &graphql.Field{Type: pageInfoType},
	},
})

var grantListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "GrantList",
	Fields: graphql.Fields{
		"results":  &graphql.Field{Type: graphql.NewList(grantType)},
		"pageInfo": &graphql.Field{Type: pageInfoType},
	},
})

var publicationListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "publicationList",
	Fields: graphql.Fields{
		"results":  &graphql.Field{Type: graphql.NewList(publicationType)},
		"pageInfo": &graphql.Field{Type: pageInfoType},
	},
})
