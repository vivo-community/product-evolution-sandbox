package graphql

//https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
import (
	"log"

	ge "github.com/OIT-ads-web/graphql_endpoint"
	"github.com/OIT-ads-web/graphql_endpoint/elastic"
	"github.com/davecgh/go-spew/spew"
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
type PagingFilter struct {
	Limit  int
	Offset int
}

type PersonFilterParam struct {
	Filter PagingFilter
}

func asPeopleFilter(params graphql.ResolveParams) (PersonFilterParam, error) {
	// default values?
	result := PersonFilterParam{PagingFilter{0, 100}}
	err := ms.Decode(params.Args, &result)
	return result, err
}

func peopleResolver(params graphql.ResolveParams) (interface{}, error) {
	// TODO: not finding a good way to default these
	// values - default is defined in graphql.InputObject
	// but then once again dealt with here
	limit := 100
	offset := 0
	// q := "*:*"
	filter, err := asPeopleFilter(params)
	if err != nil {
		limit = filter.Filter.Limit
		offset = filter.Filter.Offset
	}

	spew.Printf("limit=%d, offset=%d\n")
	personList, err := elastic.FindPeople(limit, offset)
	return personList, err
}

func publicationResolver(params graphql.ResolveParams) (interface{}, error) {
	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	publications, err := elastic.FindPublications(size, from)
	return publications, err
}

func personPublicationResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(ge.Person)

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	publicationList, err := elastic.FindPersonPublications(person.Id, size, from)
	return func() (interface{}, error) {
		return &publicationList, err
	}, nil
}

func grantResolver(params graphql.ResolveParams) (interface{}, error) {
	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	grants, err := elastic.FindGrants(size, from)
	return grants, err
}

func personGrantResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(ge.Person)

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	grants, err := elastic.FindPersonGrants(person.Id, size, from)

	return func() (interface{}, error) {
		return &grants, err
	}, nil
}
