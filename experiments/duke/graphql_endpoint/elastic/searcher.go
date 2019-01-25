package elastic

//https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
import (
	"context"
	"encoding/json"
	"fmt"
	ge "github.com/OIT-ads-web/graphql_endpoint"
	"github.com/davecgh/go-spew/spew"
	"github.com/olivere/elastic"
	"log"
)

// NOTE: these should all take a 'context' parameter
func FindPerson(personId string) (ge.Person, error) {
	var person = ge.Person{}

	ctx := context.Background()
	client := GetClient()

	log.Printf("looking for person %s\n", personId)

	get1, err := client.Get().
		Index("people").
		Id(personId).
		Do(ctx)
	if err != nil {
		return person, err
	}

	err = json.Unmarshal(*get1.Source, &person)
	return person, err
}

func FindPeople(from int, size int) (ge.PersonList, error) {
	var people []ge.Person
	ctx := context.Background()
	client := GetClient()

	q := elastic.NewMatchAllQuery()

	log.Println("looking for people")

	searchResult, err := client.Search().
		Index("people").
		Query(q).
		From(from).
		Size(size).
		// Timeout("1000ms"). or
		// Timeout(1000).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	for _, hit := range searchResult.Hits.Hits {
		person := ge.Person{}
		err := json.Unmarshal(*hit.Source, &person)
		if err != nil {
			panic(err)
		}
		people = append(people, person)
	}

	totalHits := int(searchResult.TotalHits())
	log.Printf("total hits: %d\n", totalHits)

	pageInfo := ge.FigurePaging(from, size, totalHits)
	personList := ge.PersonList{Results: people, PageInfo: pageInfo}
	return personList, err
}

func FindPublications(from int, size int) ([]ge.Publication, error) {
	var publications []ge.Publication
	ctx := context.Background()
	// should query elastic here
	client := GetClient()

	q := elastic.NewMatchAllQuery()

	searchResult, err := client.Search().
		Index("publications").
		Query(q).
		From(from).
		Size(size).
		//Pretty(true).
		// Timeout("1000ms"). or
		// Timeout(1000).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	for _, hit := range searchResult.Hits.Hits {
		publication := ge.Publication{}
		err := json.Unmarshal(*hit.Source, &publication)
		if err != nil {
			panic(err)
		}
		publications = append(publications, publication)
	}
	return publications, err
}

func FindPersonPublications(personId string, from int, size int) (ge.PublicationList, error) {
	var publications []ge.Publication
	var publicationIds []string

	ctx := context.Background()
	client := GetClient()

	q := elastic.NewMatchQuery("personId", personId)

	searchResult, err := client.Search().
		Index("authorships").
		Query(q).
		From(from).
		Size(size).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	// FIXME: could optimize better - dataloader etc...
	for _, hit := range searchResult.Hits.Hits {
		authorship := ge.Authorship{}
		err := json.Unmarshal(*hit.Source, &authorship)
		if err != nil {
			panic(err)
		}

		publicationId := authorship.PublicationId
		publicationIds = append(publicationIds, publicationId)
	}

	// NOTE: need to have the count be authorship search
	// not publication search - since pub search is just
	// an id search derived from authorship search
	totalHits := int(searchResult.TotalHits())
	log.Printf("total authorships: %d\n", totalHits)

	// ids query
	pubQuery := elastic.NewIdsQuery("publication").
		Ids(publicationIds...)

	pubResults, err := client.Search().
		Index("publications").
		Query(pubQuery).
		RequestCache(true).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	for _, hit := range pubResults.Hits.Hits {
		publication := ge.Publication{}
		err := json.Unmarshal(*hit.Source, &publication)
		if err != nil {
			panic(err)
		}
		publications = append(publications, publication)
	}

	pageInfo := ge.FigurePaging(from, size, totalHits)
	publicationList := ge.PublicationList{Results: publications, PageInfo: pageInfo}

	return publicationList, err
}

func FindGrants(personId string, from int, size int) ([]ge.Grant, error) {
	var grants []ge.Grant
	var grantIds []string

	ctx := context.Background()
	client := GetClient()

	q := elastic.NewMatchQuery("personId", personId)

	searchResult, err := client.Search().
		Index("funding-roles").
		Query(q).
		From(from).
		Size(size).
		Do(ctx)
	if err != nil {
		// handle error
		panic(err)
	}

	// fixme: could optimize better - dataloader etc...
	for _, hit := range searchResult.Hits.Hits {
		fundingRole := ge.FundingRole{}
		err := json.Unmarshal(*hit.Source, &fundingRole)
		if err != nil {
			panic(err)
		}

		grantId := fundingRole.GrantId
		grantIds = append(grantIds, grantId)
	}

	grantQuery := elastic.NewIdsQuery("grant").
		Ids(grantIds...)

	grantResults, err := client.Search().
		Index("grants").
		Query(grantQuery).
		RequestCache(true).
		Do(ctx)
	if err != nil {
		// handle error
		panic(err)
	}
	for _, hit := range grantResults.Hits.Hits {
		grant := ge.Grant{}
		err := json.Unmarshal(*hit.Source, &grant)
		if err != nil {
			panic(err)
		}
		grants = append(grants, grant)
	}

	return grants, err
}

func FindAffiliations(personId string, from int, size int) ([]ge.Affiliation, error) {
	var affiliations []ge.Affiliation

	ctx := context.Background()
	client := GetClient()

	q := elastic.NewMatchQuery("personId", personId)

	searchResult, err := client.Search().
		Index("affiliations").
		Query(q).
		From(from).
		Size(size).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	for _, hit := range searchResult.Hits.Hits {
		affiliation := ge.Affiliation{}
		err := json.Unmarshal(*hit.Source, &affiliation)
		if err != nil {
			panic(err)
		}
		affiliations = append(affiliations, affiliation)

	}
	return affiliations, err
}

func FindEducations(personId string, from int, size int) ([]ge.Education, error) {
	var educations []ge.Education

	ctx := context.Background()
	client := GetClient()

	q := elastic.NewMatchQuery("personId", personId)

	searchResult, err := client.Search().
		Index("educations").
		Query(q).
		From(from).
		Size(size).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	for _, hit := range searchResult.Hits.Hits {
		education := ge.Education{}
		err := json.Unmarshal(*hit.Source, &education)
		if err != nil {
			panic(err)
		}
		educations = append(educations, education)
	}
	return educations, err
}

func ListAll(index string) {
	ctx := context.Background()
	client := GetClient()
	q := elastic.NewMatchAllQuery()

	searchResult, err := client.Search().
		Index(index).
		//Type().
		Query(q).
		From(100).
		Size(100).
		Pretty(true).
		// Timeout("1000ms"). or
		// Timeout(1000).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
		
	fmt.Println("********* BEGIN **********")
	for _, hit := range searchResult.Hits.Hits {
		var obj interface{}
		err := json.Unmarshal(*hit.Source, &obj)
		if err != nil {
			panic(err)
		}
		spew.Println(obj)
	}
	fmt.Printf("********* END (%d) **********\n", searchResult.TotalHits())
}

func IdQuery(index string, ids []string) {
	// NOTE: can send 'type' into NewIdsQuery
	q := elastic.NewIdsQuery().Ids(ids...) //.QueryName("my_query")
	ctx := context.Background()
	client := GetClient()

	searchResult, err := client.Search().
		Index(index).
		//Type().
		Query(q).
		From(0).
		Size(1000).
		Pretty(true).
		// Timeout("1000ms"). or
		// Timeout(1000).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Println("********* BEGIN **********")
	for _, hit := range searchResult.Hits.Hits {
		var obj interface{}
		err := json.Unmarshal(*hit.Source, &obj)
		if err != nil {
			panic(err)
		}
		spew.Println(obj)
	}
    fmt.Println("************** END **********")
}

func FindOne(index string, personId string) {
	ctx := context.Background()
	client := GetClient()

	get1, err := client.Get().
		Index(index).
		Id(personId).
		Do(ctx)

	switch {
	case elastic.IsNotFound(err):
		fmt.Println("404 not found")
	case elastic.IsConnErr(err):
		fmt.Println("connectino error")
	case elastic.IsTimeout(err):
		fmt.Println("timeout")
	case err != nil:
		panic(err)
	}

    var obj interface{}
    err = json.Unmarshal(*get1.Source, &obj)
	if err != nil {
		panic(err)
	}
	spew.Println(obj)
}
