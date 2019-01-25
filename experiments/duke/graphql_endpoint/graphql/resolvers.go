package graphql

//https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
import (
	ge "github.com/OIT-ads-web/graphql_endpoint"
	"github.com/OIT-ads-web/graphql_endpoint/elastic"
	"github.com/graphql-go/graphql"
	"log"
)

func personResolver(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)
	log.Printf("looking for person %s\n", id)

	person, err := elastic.FindPerson(id)
	return person, err
}

func peopleResolver(params graphql.ResolveParams) (interface{}, error) {
	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	personList, err := elastic.FindPeople(size, from)
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

	publicationList, err := elastic.FindPersonPublications(person.Id, from, size)
	return func() (interface{}, error) {
		return &publicationList, err
	}, nil
}

func grantResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(ge.Person)

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	grants, err := elastic.FindGrants(person.Id, from, size)

	return func() (interface{}, error) {
		return &grants, err
	}, nil
}

func affiliationResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(ge.Person)
	var affiliations []ge.Affiliation

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	affiliations, err := elastic.FindAffiliations(person.Id, size, from)
	return func() (interface{}, error) {
		return &affiliations, err
	}, nil
}

func educationResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(ge.Person)

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	educations, err := elastic.FindEducations(person.Id, from, size)
	return func() (interface{}, error) {
		return &educations, err
	}, nil
}
