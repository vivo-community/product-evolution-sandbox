package elastic

//https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	ge "github.com/OIT-ads-web/graphql_endpoint"
	"github.com/davecgh/go-spew/spew"
	"github.com/olivere/elastic"
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

func FindPublication(publicationId string) (ge.Publication, error) {
	var publication = ge.Publication{}

	ctx := context.Background()
	client := GetClient()

	log.Printf("looking for publication %s\n", publicationId)

	get1, err := client.Get().
		Index("publications").
		Id(publicationId).
		Do(ctx)
	if err != nil {
		return publication, err
	}

	err = json.Unmarshal(*get1.Source, &publication)
	return publication, err
}

func FindGrant(grantId string) (ge.Grant, error) {
	var grant = ge.Grant{}

	ctx := context.Background()
	client := GetClient()

	log.Printf("looking for grant %s\n", grantId)

	get1, err := client.Get().
		Index("grants").
		Id(grantId).
		Do(ctx)
	if err != nil {
		return grant, err
	}

	err = json.Unmarshal(*get1.Source, &grant)
	return grant, err
}

func parsePeopleAggregations(facets elastic.Aggregations) *ge.PeopleFacets {
	peopleFacets := &ge.PeopleFacets{}

	if agg, found := facets.Nested("keywords"); found {
		var facets []ge.Facet
		if sub, subFound := agg.Terms("keyword"); subFound {
			for _, bucket := range sub.Buckets {
				facet := ge.Facet{Label: bucket.Key.(string), Count: bucket.DocCount}
				facets = append(facets, facet)
			}
		}
		peopleFacets.Keywords = facets
	}

	if agg, found := facets.Nested("affiliations"); found {
		var facets []ge.Facet
		if sub, subFound := agg.Terms("department"); subFound {
			for _, bucket := range sub.Buckets {
				facet := ge.Facet{Label: bucket.Key.(string), Count: bucket.DocCount}
				facets = append(facets, facet)
			}
		}
		peopleFacets.Departments = facets
	}

	if agg, found := facets.Terms("types"); found {
		var facets []ge.Facet
		for _, bucket := range agg.Buckets {
			facet := ge.Facet{Label: bucket.Key.(string), Count: bucket.DocCount}
			facets = append(facets, facet)
		}
		peopleFacets.Types = facets
	}
	return peopleFacets
}

func FindPeople(limit int, offset int, query string) (ge.PersonList, error) {
	var people []ge.Person
	ctx := context.Background()
	client := GetClient()

	q := elastic.NewQueryStringQuery(query)
	log.Println("looking for people")

	service := client.Search().
		Index("people").
		Query(q).
		From(offset).
		Size(limit)

	keywordsSize := 100
	departmentsSize := 100

	// TODO: kind of kludged together here - probably much better way to do
	agg := elastic.NewTermsAggregation().Field("type.label")
	service = service.Aggregation("types", agg)

	nested := elastic.NewNestedAggregation().Path("keywordList")
	subAgg := nested.SubAggregation("keyword",
		elastic.NewTermsAggregation().
			Field("keywordList.label.keyword").
			Size(keywordsSize))

	nested2 := elastic.NewNestedAggregation().Path("affiliationList")
	subAgg2 := nested2.SubAggregation("department",
		elastic.NewTermsAggregation().
			Field("affiliationList.organization.label.dept").
			Size(departmentsSize))

	service = service.Aggregation("keywords", subAgg)
	service = service.Aggregation("affiliations", subAgg2)

	searchResult, err := service.Do(ctx)
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

	// TODO: might be one off
	totalHits := int(searchResult.TotalHits())
	log.Printf("total hits: %d\n", totalHits)

	pageInfo := ge.FigurePaging(limit, offset, totalHits)
	facets := parsePeopleAggregations(searchResult.Aggregations)
	personList := ge.PersonList{Results: people, PageInfo: pageInfo, Facets: facets}
	return personList, err
}

func FindPublications(limit int, offset int, query string) (ge.PublicationList, error) {
	var publications []ge.Publication
	ctx := context.Background()
	// should query elastic here
	client := GetClient()

	//q := elastic.NewMatchAllQuery()
	q := elastic.NewQueryStringQuery(query)
	service := client.Search().
		Index("publications").
		Query(q).
		From(offset).
		Size(limit)

	/*
		// TODO: kind of kludged together here - probably much better way to do
		agg := elastic.NewTermsAggregation().Field("type.label")
		service = service.Aggregation("types", agg)

		nested := elastic.NewNestedAggregation().Path("keywordList")
		subAgg := nested.SubAggregation("keyword", elastic.NewTermsAggregation().Field("keywordList.label.keyword"))

		nested2 := elastic.NewNestedAggregation().Path("affiliationList")
		subAgg2 := nested2.SubAggregation("department",
			elastic.NewTermsAggregation().Field("affiliationList.organization.label.dept"))

		service = service.Aggregation("keywords", subAgg)
		service = service.Aggregation("affiliations", subAgg2)
	*/

	searchResult, err := service.Do(ctx)
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

	// might be one off
	totalHits := int(searchResult.TotalHits())
	log.Printf("total hits: %d\n", totalHits)

	pageInfo := ge.FigurePaging(limit, offset, totalHits)
	// eventually
	//facets := parsePublicationsAggregations(searchResult.Aggregations)
	//publicationList := ge.PublicationList{Results: publications, PageInfo: pageInfo, Facets: facets}
	publicationList := ge.PublicationList{Results: publications, PageInfo: pageInfo}
	return publicationList, err
	//return publications, err
}

func FindPersonPublications(personId string, limit int, offset int) (ge.PublicationList, error) {
	var publications []ge.Publication
	var publicationIds []string

	ctx := context.Background()
	client := GetClient()

	q := elastic.NewMatchQuery("personId", personId)

	searchResult, err := client.Search().
		Index("authorships").
		Query(q).
		From(offset).
		Size(limit).
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

	// NOTE: not sure this is actually correct, might be one off
	totalHits := int(searchResult.TotalHits())
	log.Printf("total authorships: %d\n", totalHits)

	// ids query
	pubQuery := elastic.NewIdsQuery("publication").
		Ids(publicationIds...)

	pubResults, err := client.Search().
		Index("publications").
		Query(pubQuery).
		From(0).
		Size(totalHits).
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

	log.Printf("size: %d, from:%d\n", limit, offset)

	pageInfo := ge.FigurePaging(limit, offset, totalHits)
	publicationList := ge.PublicationList{Results: publications, PageInfo: pageInfo}

	return publicationList, err
}

func FindGrants(limit int, offset int, query string) (ge.GrantList, error) {
	var grants []ge.Grant
	ctx := context.Background()
	client := GetClient()

	//q := elastic.NewMatchAllQuery()
	q := elastic.NewQueryStringQuery(query)
	service := client.Search().
		Index("grants").
		Query(q).
		From(limit).
		Size(offset)

	searchResult, err := service.Do(ctx)

	if err != nil {
		// Handle error
		panic(err)
	}

	for _, hit := range searchResult.Hits.Hits {
		grant := ge.Grant{}
		err := json.Unmarshal(*hit.Source, &grant)
		if err != nil {
			panic(err)
		}
		grants = append(grants, grant)
	}

	// is this the correct number?
	totalHits := int(searchResult.TotalHits())
	log.Printf("total hits: %d\n", totalHits)

	pageInfo := ge.FigurePaging(limit, offset, totalHits)
	// eventually
	//facets := parseGrantsAggregations(searchResult.Aggregations)
	//grantList := ge.GrantList{Results: grants, PageInfo: pageInfo, Facets: facets}
	grantList := ge.GrantList{Results: grants, PageInfo: pageInfo}
	return grantList, err
}

func FindPersonGrants(personId string, limit int, offset int) (ge.GrantList, error) {
	var grants []ge.Grant
	var grantIds []string

	ctx := context.Background()
	client := GetClient()

	q := elastic.NewMatchQuery("personId", personId)

	searchResult, err := client.Search().
		Index("funding-roles").
		Query(q).
		From(offset).
		Size(limit).
		Do(ctx)
	if err != nil {
		// handle error
		panic(err)
	}

	// is this the correct number?
	totalHits := int(searchResult.TotalHits())
	log.Printf("total funding-roles: %d\n", totalHits)

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
		From(0).
		Size(totalHits).
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

	pageInfo := ge.FigurePaging(limit, offset, totalHits)
	grantList := ge.GrantList{Results: grants, PageInfo: pageInfo}
	return grantList, err
}

// remaining are just debug/util functions
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
		spew.Printf("%v\n", obj)
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
		spew.Printf("%v\n", obj)
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
	spew.Printf("%v\n", obj)
}
