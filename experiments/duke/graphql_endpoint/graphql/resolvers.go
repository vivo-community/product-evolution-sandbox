package graphql

//https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
import (
	"log"

	ge "github.com/OIT-ads-web/graphql_endpoint"
	"github.com/OIT-ads-web/graphql_endpoint/elastic"
	"github.com/graphql-go/graphql"
)

func personResolver(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)
	log.Printf("looking for person %s\n", id)

	person, err := elastic.FindPerson(id)
	return person, err
}

func peopleResolver(params graphql.ResolveParams) (interface{}, error) {
	// TODO: how else to get default?
	limit := 100
	offset := 0
	if filter, ok := params.Args["filter"]; ok {
		obj := filter.(map[string]interface{})
		limit = obj["limit"].(int)
		offset = obj["offset"].(int)
	}

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
