package graphql

//https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
import (
	"fmt"
	"log"

	ge "github.com/OIT-ads-web/graphql_endpoint"
	"github.com/OIT-ads-web/graphql_endpoint/elastic"
	"github.com/graphql-go/graphql"

	ms "github.com/mitchellh/mapstructure"
)

func personResolver(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)
	log.Printf("looking for person %s\n", id)

	person, err := elastic.FindPerson(id)
	return person, err
}

// NOTE: this duplicates structure here:
// var PersonFilter *graphql.InputObject
// not sure best way to go about this
type CommonFilter struct {
	Limit  int
	Offset int
	Query  string
}

// NOTE: these aren't different now, but dealing with
// facets would probably make them different
//  `mapstructure:",squash"`
type PersonFilterParam struct {
	Filter CommonFilter
}

type PublicationFilterParam struct {
	Filter CommonFilter
}

type GrantFilterParam struct {
	Filter CommonFilter
}

func convertPeopleFilter(params graphql.ResolveParams) (PersonFilterParam, error) {
	result := PersonFilterParam{}
	err := ms.Decode(params.Args, &result)
	return result, err
}

func convertPublicationFilter(params graphql.ResolveParams) (PublicationFilterParam, error) {
	// default values?
	result := PublicationFilterParam{}
	err := ms.Decode(params.Args, &result)
	return result, err
}

func convertGrantFilter(params graphql.ResolveParams) (GrantFilterParam, error) {
	// default values?
	result := GrantFilterParam{}
	err := ms.Decode(params.Args, &result)
	return result, err
}

func peopleResolver(params graphql.ResolveParams) (interface{}, error) {
	// TODO: not finding a good way to default these
	// e.g. if filter is not sent at all
	limit := 100
	offset := 0
	query := "*:*"
	filter, err := convertPeopleFilter(params)

	if err == nil {
		limit = filter.Filter.Limit
		offset = filter.Filter.Offset
		// NOTE: this is not that great
		query = fmt.Sprintf("*:%v*", filter.Filter.Query)
	}

	personList, err := elastic.FindPeople(limit, offset, query)
	return personList, err
}

func publicationResolver(params graphql.ResolveParams) (interface{}, error) {
	// TODO: not finding a good way to default these
	limit := 100
	offset := 0
	query := "*:*"
	filter, err := convertPublicationFilter(params)
	if err == nil {
		limit = filter.Filter.Limit
		offset = filter.Filter.Offset
		query = fmt.Sprintf("*:%v*", filter.Filter.Query)
	}

	publications, err := elastic.FindPublications(limit, offset, query)
	return publications, err
}

func personPublicationResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(ge.Person)

	limit := params.Args["limit"].(int)
	offset := params.Args["offset"].(int)

	publicationList, err := elastic.FindPersonPublications(person.Id, limit, offset)
	return func() (interface{}, error) {
		return &publicationList, err
	}, nil
}

func grantResolver(params graphql.ResolveParams) (interface{}, error) {
	limit := 100
	offset := 0
	query := "*:*"
	filter, err := convertGrantFilter(params)
	if err == nil {
		limit = filter.Filter.Limit
		offset = filter.Filter.Offset
		query = fmt.Sprintf("*:%v*", filter.Filter.Query)
	}
	grants, err := elastic.FindGrants(limit, offset, query)
	return grants, err
}

func personGrantResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(ge.Person)

	limit := params.Args["limit"].(int)
	offset := params.Args["offset"].(int)

	grants, err := elastic.FindPersonGrants(person.Id, limit, offset)

	return func() (interface{}, error) {
		return &grants, err
	}, nil
}
