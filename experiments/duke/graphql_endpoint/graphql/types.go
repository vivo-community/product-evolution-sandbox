package graphql

import (
	"github.com/graphql-go/graphql"
)

// NOTE: removed Type from end of names per
// https://dave.cheney.net/practical-go/presentations/qcon-china.html
var pageInfo = graphql.NewObject(graphql.ObjectConfig{
	Name: "PageInfo",
	Fields: graphql.Fields{
		"perPage":    &graphql.Field{Type: graphql.Int},
		"page":       &graphql.Field{Type: graphql.Int},
		"totalPages": &graphql.Field{Type: graphql.Int},
		"count":      &graphql.Field{Type: graphql.Int},
	},
})

var grant = graphql.NewObject(graphql.ObjectConfig{
	Name: "Grant",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"label":     &graphql.Field{Type: graphql.String},
		"roleName":  &graphql.Field{Type: graphql.String},
		"startDate": &graphql.Field{Type: dateResolution},
		"endDate":   &graphql.Field{Type: dateResolution},
	},
})

var organization = graphql.NewObject(graphql.ObjectConfig{
	Name: "Organization",
	Fields: graphql.Fields{
		"id":    &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
	},
})

var education = graphql.NewObject(graphql.ObjectConfig{
	Name: "Education",
	Fields: graphql.Fields{
		"credential":             &graphql.Field{Type: graphql.String},
		"credentialAbbreviation": &graphql.Field{Type: graphql.String},
		"organization":           &graphql.Field{Type: organization},
		"dateReceived":           &graphql.Field{Type: dateResolution},
	},
})

var affiliation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Affiliation",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"label":     &graphql.Field{Type: graphql.String},
		"startDate": &graphql.Field{Type: dateResolution},
	},
})

var keyword = graphql.NewObject(graphql.ObjectConfig{
	Name: "Keyword",
	Fields: graphql.Fields{
		"uri":   &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
	},
})

var extension = graphql.NewObject(graphql.ObjectConfig{
	Name: "Extension",
	Fields: graphql.Fields{
		"key":   &graphql.Field{Type: graphql.String},
		"value": &graphql.Field{Type: graphql.String},
	},
})

var dateResolution = graphql.NewObject(graphql.ObjectConfig{
	Name: "DateResolution",
	Fields: graphql.Fields{
		"dateTime":   &graphql.Field{Type: graphql.String},
		"resolution": &graphql.Field{Type: graphql.String},
	},
})

var personName = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonName",
	Fields: graphql.Fields{
		"firstName":  &graphql.Field{Type: graphql.String},
		"lastName":   &graphql.Field{Type: graphql.String},
		"middleName": &graphql.Field{Type: graphql.String},
	},
})

var personImage = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonImage",
	Fields: graphql.Fields{
		"main":      &graphql.Field{Type: graphql.String},
		"thumbnail": &graphql.Field{Type: graphql.String},
	},
})

var typeOf = graphql.NewObject(graphql.ObjectConfig{
	Name: "Type",
	Fields: graphql.Fields{
		"code":  &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
	},
})

var personIdentifier = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonIdentifier",
	Fields: graphql.Fields{
		"orchid": &graphql.Field{Type: graphql.String},
		"isni":   &graphql.Field{Type: graphql.String},
	},
})

var personContact = graphql.NewObject(graphql.ObjectConfig{
	Name: "Contact",
	Fields: graphql.Fields{
		"emailList":    &graphql.Field{Type: graphql.NewList(email)},
		"locationList": &graphql.Field{Type: graphql.NewList(location)},
		"phoneList":    &graphql.Field{Type: graphql.NewList(phone)},
		"websiteList":  &graphql.Field{Type: graphql.NewList(website)},
	},
})

var serviceRole = graphql.NewObject(graphql.ObjectConfig{
	Name: "ServiceRole",
	Fields: graphql.Fields{
		"uri":          &graphql.Field{Type: graphql.String},
		"label":        &graphql.Field{Type: graphql.String},
		"description":  &graphql.Field{Type: graphql.String},
		"startDate":    &graphql.Field{Type: dateResolution},
		"endDate":      &graphql.Field{Type: dateResolution},
		"organization": &graphql.Field{Type: organization},
		"type":         &graphql.Field{Type: typeOf},
	},
})

var email = graphql.NewObject(graphql.ObjectConfig{
	Name: "Email",
	Fields: graphql.Fields{
		"uri":   &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
		"type":  &graphql.Field{Type: typeOf},
	},
})

var phone = graphql.NewObject(graphql.ObjectConfig{
	Name: "Phone",
	Fields: graphql.Fields{
		"uri":   &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
		"type":  &graphql.Field{Type: typeOf},
	},
})

var location = graphql.NewObject(graphql.ObjectConfig{
	Name: "Location",
	Fields: graphql.Fields{
		"uri":   &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
		"type":  &graphql.Field{Type: typeOf},
	},
})

var website = graphql.NewObject(graphql.ObjectConfig{
	Name: "Website",
	Fields: graphql.Fields{
		"uri":   &graphql.Field{Type: graphql.String},
		"url":   &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
		"type":  &graphql.Field{Type: typeOf},
	},
})

var courseTaught = graphql.NewObject(graphql.ObjectConfig{
	Name: "CourseTaught",
	Fields: graphql.Fields{
		"uri":          &graphql.Field{Type: graphql.String},
		"subject":      &graphql.Field{Type: graphql.String},
		"role":         &graphql.Field{Type: graphql.String},
		"courseName":   &graphql.Field{Type: graphql.String},
		"courseNumber": &graphql.Field{Type: graphql.String},
		"startDate":    &graphql.Field{Type: dateResolution},
		"endDate":      &graphql.Field{Type: dateResolution},
		"organization": &graphql.Field{Type: organization},
		"type":         &graphql.Field{Type: typeOf},
	},
})

var publicationIdentifier = graphql.NewObject(graphql.ObjectConfig{
	Name: "PublicationIdentifier",
	Fields: graphql.Fields{
		"isbn10": &graphql.Field{Type: graphql.String},
		"isbn13": &graphql.Field{Type: graphql.String},
		"pmid":   &graphql.Field{Type: graphql.String},
		"doi":    &graphql.Field{Type: graphql.String},
		"pmcid":  &graphql.Field{Type: graphql.String},
	},
})

var publicationVenue = graphql.NewObject(graphql.ObjectConfig{
	Name: "PublicationVenue",
	Fields: graphql.Fields{
		"uri":   &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
	},
})

var publication = graphql.NewObject(graphql.ObjectConfig{
	Name: "Publication",
	Fields: graphql.Fields{
		"id":               &graphql.Field{Type: graphql.String},
		"uri":              &graphql.Field{Type: graphql.String},
		"title":            &graphql.Field{Type: graphql.String},
		"identifier":       &graphql.Field{Type: publicationIdentifier},
		"type":             &graphql.Field{Type: typeOf},
		"authorList":       &graphql.Field{Type: graphql.String},
		"venue":            &graphql.Field{Type: publicationVenue},
		"abstract":         &graphql.Field{Type: graphql.String},
		"dateStandardized": &graphql.Field{Type: dateResolution},
		"dateDisplay":      &graphql.Field{Type: graphql.String},
		"pageRange":        &graphql.Field{Type: graphql.String},
		"pageStart":        &graphql.Field{Type: graphql.String},
		"pageEnd":          &graphql.Field{Type: graphql.String},
		"issue":            &graphql.Field{Type: graphql.String},
		"volume":           &graphql.Field{Type: graphql.String},
	},
})

var overview = graphql.NewObject(graphql.ObjectConfig{
	Name: "Overview",
	Fields: graphql.Fields{
		"code":  &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
	},
})

var authorship = graphql.NewObject(graphql.ObjectConfig{
	Name: "Authorship",
	Fields: graphql.Fields{
		"id":            &graphql.Field{Type: graphql.String},
		"uri":           &graphql.Field{Type: graphql.String},
		"publicationId": &graphql.Field{Type: graphql.String},
		"personId":      &graphql.Field{Type: graphql.String},
		"label":         &graphql.Field{Type: graphql.String},
	},
})

var person = graphql.NewObject(graphql.ObjectConfig{
	Name: "Person",
	Fields: graphql.Fields{
		"uri":              &graphql.Field{Type: graphql.String},
		"id":               &graphql.Field{Type: graphql.String},
		"sourceId":         &graphql.Field{Type: graphql.String},
		"primaryTitle":     &graphql.Field{Type: graphql.String},
		"name":             &graphql.Field{Type: personName},
		"image":            &graphql.Field{Type: personImage},
		"identifier":       &graphql.Field{Type: personIdentifier},
		"contact":          &graphql.Field{Type: personContact},
		"type":             &graphql.Field{Type: typeOf},
		"overviewList":     &graphql.Field{Type: graphql.NewList(overview)},
		"keywordList":      &graphql.Field{Type: graphql.NewList(keyword)},
		"extensionList":    &graphql.Field{Type: graphql.NewList(extension)},
		"affliationList":   &graphql.Field{Type: graphql.NewList(affiliation)},
		"educationList":    &graphql.Field{Type: graphql.NewList(education)},
		"serviceRoleList":  &graphql.Field{Type: graphql.NewList(serviceRole)},
		"courseTaughtList": &graphql.Field{Type: graphql.NewList(courseTaught)},
		// these can be paged, since they involve further queries
		"publicationList": &graphql.Field{
			Type: publicationList,
			Args: graphql.FieldConfigArgument{
				"size": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 100},
				"from": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: personPublicationResolver,
		},
		"grantList": &graphql.Field{
			Type: grantList,
			Args: graphql.FieldConfigArgument{
				"size": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 100},
				"from": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: personGrantResolver,
		},
	},
})


/*
FIXME: not sure best way to send in sorting paramters right now
*/
var sortField = graphql.NewObject(graphql.ObjectConfig{
	Name: "SortField",
	Fields: graphql.Fields{
		"name":  &graphql.Field{Type: graphql.String},
		"order": &graphql.Field{Type: graphql.String},
		"mode":  &graphql.Field{Type: graphql.String},
	},
})

var sorter = graphql.NewObject(graphql.ObjectConfig{
	Name: "Sort",
	Fields: graphql.Fields{
		"fields": &graphql.Field{Type: graphql.NewList(sortField)},
	},
})

var filter = graphql.NewObject(graphql.ObjectConfig{
	Name: "Filter",
	Fields: graphql.Fields{
		"limit":  &graphql.Field{Type: graphql.Int},
		"offset": &graphql.Field{Type: graphql.Int},
		"sort":   &graphql.Field{Type: sorter},
	},
})

var personFacets = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonFacets",
	Fields: graphql.Fields{
		"departments": &graphql.Field{Type: graphql.NewList(facet)},
		"types":       &graphql.Field{Type: graphql.NewList(facet)},
		"keywords":    &graphql.Field{Type: graphql.NewList(facet)},
	},
})

var facet = graphql.NewObject(graphql.ObjectConfig{
	Name: "Facet",
	Fields: graphql.Fields{
		"label": &graphql.Field{Type: graphql.String},
		"count": &graphql.Field{Type: graphql.Int},
	},
})

var personList = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonList",
	Fields: graphql.Fields{
		"results":  &graphql.Field{Type: graphql.NewList(person)},
		"pageInfo": &graphql.Field{Type: pageInfo},
		"facets":   &graphql.Field{Type: personFacets},
	},
})

var grantList = graphql.NewObject(graphql.ObjectConfig{
	Name: "GrantList",
	Fields: graphql.Fields{
		"results":  &graphql.Field{Type: graphql.NewList(grant)},
		"pageInfo": &graphql.Field{Type: pageInfo},
	},
})

var publicationList = graphql.NewObject(graphql.ObjectConfig{
	Name: "publicationList",
	Fields: graphql.Fields{
		"results":  &graphql.Field{Type: graphql.NewList(publication)},
		"pageInfo": &graphql.Field{Type: pageInfo},
	},
})
