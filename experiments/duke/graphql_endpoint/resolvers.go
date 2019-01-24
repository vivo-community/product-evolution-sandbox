package graphql_endpoint

import (
	"context"
	"encoding/json"
	"github.com/graphql-go/graphql"
	"github.com/olivere/elastic"
	"log"
)

func FigurePaging(from int, size int, totalHits int) PageInfo {
	// has to at least be page 1, maybe even if totalHits = 0
	var currentPage = 1
	var offset = from

	if (offset / size) > 0 {
		if (offset % size) > 0 {
			currentPage = (offset / size) + 1
		} else {
			currentPage = (offset / size) - 1
		}
	}
	var totalPages = totalHits / size
	var remainder = totalHits % size
	if remainder > 0 {
		totalPages += 1
	}
	pageInfo := PageInfo{PerPage: size,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		Count:       totalHits}
	return pageInfo
}

func personResolver(params graphql.ResolveParams) (interface{}, error) {
	var person = Person{}
	ctx := context.Background()
	client := GetClient()

	id := params.Args["id"].(string)
	log.Printf("looking for person %s\n", id)

	get1, err := client.Get().
		Index("people").
		Id(id).
		Do(ctx)
	if err != nil {
		return person, err
	}

	err = json.Unmarshal(*get1.Source, &person)

	if err != nil {
		return person, err
	}
	return person, nil
}

func peopleResolver(params graphql.ResolveParams) (interface{}, error) {
	var people []Person
	ctx := context.Background()
	client := GetClient()

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

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
		person := Person{}
		err := json.Unmarshal(*hit.Source, &person)
		if err != nil {
			panic(err)
		}
		people = append(people, person)
	}

	totalHits := int(searchResult.TotalHits())
	log.Printf("total hits: %d\n", totalHits)

	pageInfo := FigurePaging(from, size, totalHits)
	personList := PersonList{Results: people, PageInfo: pageInfo}
	return personList, nil
}

func publicationResolver(params graphql.ResolveParams) (interface{}, error) {
	var publications []Publication
	ctx := context.Background()
	// should query elastic here
	client := GetClient()

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

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
		publication := Publication{}
		err := json.Unmarshal(*hit.Source, &publication)
		if err != nil {
			panic(err)
		}
		publications = append(publications, publication)
	}
	return publications, nil
}

func personPublicationResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(Person)
	var publications []Publication
	var publicationIds []string

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	ctx := context.Background()
	client := GetClient()

	q := elastic.NewMatchQuery("personId", person.Id)

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
		authorship := Authorship{}
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
		publication := Publication{}
		err := json.Unmarshal(*hit.Source, &publication)
		if err != nil {
			panic(err)
		}
		publications = append(publications, publication)
	}

	pageInfo := FigurePaging(from, size, totalHits)
	publicationList := PublicationList{Results: publications, PageInfo: pageInfo}
	return func() (interface{}, error) {
		return &publicationList, nil
	}, nil
}

//
func grantResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(Person)
	var grants []Grant
	var grantIds []string

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	ctx := context.Background()
	client := GetClient()

	q := elastic.NewMatchQuery("personId", person.Id)

	searchResult, err := client.Search().
		Index("funding-roles").
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
		fundingRole := FundingRole{}
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
		// Handle error
		panic(err)
	}
	for _, hit := range grantResults.Hits.Hits {
		grant := Grant{}
		err := json.Unmarshal(*hit.Source, &grant)
		if err != nil {
			panic(err)
		}
		grants = append(grants, grant)
	}

	return func() (interface{}, error) {
		return &grants, nil
	}, nil
	//return grants, nil
}

func affiliationResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(Person)
	var affiliations []Affiliation

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	ctx := context.Background()
	client := GetClient()

	q := elastic.NewMatchQuery("personId", person.Id)

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
		affiliation := Affiliation{}
		err := json.Unmarshal(*hit.Source, &affiliation)
		if err != nil {
			panic(err)
		}
		affiliations = append(affiliations, affiliation)

	}
	return func() (interface{}, error) {
		return &affiliations, nil
	}, nil
	//return affiliations, nil
}

func educationResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(Person)
	var educations []Education

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	ctx := context.Background()
	client := GetClient()

	q := elastic.NewMatchQuery("personId", person.Id)

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
		education := Education{}
		err := json.Unmarshal(*hit.Source, &education)
		if err != nil {
			panic(err)
		}
		educations = append(educations, education)
	}
	return func() (interface{}, error) {
		return &educations, nil
	}, nil
	//return educations, nil
}
