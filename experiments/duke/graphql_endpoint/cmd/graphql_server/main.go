package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/OIT-ads-web/graphql_endpoint"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/olivere/elastic"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

/*
type Config struct {
	Elastic elasticSearch `toml:"elastic"`
	Graphql graphqlServer `toml:"graphql"`
}

type elasticSearch struct {
	Url string
}

type graphqlServer struct {
	Port int
}

var client *elastic.Client
var conf Config

func GetClient() *elastic.Client {
	return client
}
*/

/*
var pageInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PageInfo",
	Fields: graphql.Fields{
		"perPage":    &graphql.Field{Type: graphql.Int},
		"page":       &graphql.Field{Type: graphql.Int},
		"totalPages": &graphql.Field{Type: graphql.Int},
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

var personTypeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonType",
	Fields: graphql.Fields{
		"code":  &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
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
		"authorList": &graphql.Field{Type: graphql.String},
		"doi":        &graphql.Field{Type: graphql.String},
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
*/

/*
func publicationResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(models.Person)
	var publications []models.Publication
	var publicationIds []string

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	ctx := context.Background()
	client = GetClient()

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
		authorship := models.Authorship{}
		err := json.Unmarshal(*hit.Source, &authorship)
		if err != nil {
			panic(err)
		}

		publicationId := authorship.PublicationId
		publicationIds = append(publicationIds, publicationId)
	}

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
		publication := models.Publication{}
		err := json.Unmarshal(*hit.Source, &publication)
		if err != nil {
			panic(err)
		}
		publications = append(publications, publication)
	}

	pageInfo := PageInfo{PerPage: size,
		Page:       (from / size) + 1,
		TotalPages: (int(pubResults.TotalHits()) / size) + 1}

	publicationList := PublicationList{Results: publications, PageInfo: pageInfo}
	//return personList, nil
	return func() (interface{}, error) {
		return &publicationList, nil
	}, nil
*/
	/*

		for _, hit := range pubResults.Hits.Hits {
			publication := models.Publication{}
			err := json.Unmarshal(*hit.Source, &publication)
			if err != nil {
				panic(err)
			}
			publications = append(publications, publication)
		}

		return func() (interface{}, error) {
			return &publications, nil
		}, nil
		//return publications, nil
	*/
/*
}
*/


/*
	    for _, hit := range pubResults.Hits.Hits {
		    publication := models.Publication{}
		    err := json.Unmarshal(*hit.Source, &publication)
		    if err != nil {
			    panic(err)
		    }
		    publications = append(publications, publication)
	    }

		pageInfo := PageInfo{PerPage: size,
			Page:       (from / size) + 1,
			TotalPages: (int(pubResults.TotalHits()) / size) + 1}

		publicationList := PublicationList{Results: publications, PageInfo: pageInfo}
		//return personList, nil
	    return func() (interface{}, error) {
		    return &publicationList, nil
	    }, nil
*/
/*
func grantResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(models.Person)
	var grants []models.Grant
	var grantIds []string

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	ctx := context.Background()
	client = GetClient()

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
		fundingRole := models.FundingRole{}
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
		grant := models.Grant{}
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
	person, _ := params.Source.(models.Person)
	var affiliations []models.Affiliation

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	ctx := context.Background()
	client = GetClient()

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
		affiliation := models.Affiliation{}
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
	person, _ := params.Source.(models.Person)
	var educations []models.Education

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	ctx := context.Background()
	client = GetClient()

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
		education := models.Education{}
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

func extensionResolver(params graphql.ResolveParams) (interface{}, error) {
	person, _ := params.Source.(models.Person)

	size := params.Args["size"].(int)
	from := params.Args["from"].(int)

	count := len(person.Extensions)
	extensions := make([]models.Extension, count)

	if count < size {
		//just get all
		copy(extensions, person.Extensions)
	}
	if from > count {
		// do nothing, return empty array
	}

	if from < count && size < count {
		copy(extensions, person.Extensions[from-1:size])
	}
	if from < count && size > count {
		copy(extensions, person.Extensions[from-1:count])
	}

	var totalPages int
    var page int

	if count > 0 {
		totalPages = (count / size) + 1
	} else {
		totalPages = 0
	}

	if count > 0 {
		page = (from / size ) + 1
	} else {
		page = 0
	}

	pageInfo := PageInfo{PerPage: size,
		Page:       page,
		TotalPages: totalPages}

	extensionList := ExtensionList{Results: extensions, PageInfo: pageInfo}

	return func() (interface{}, error) {
		return &extensionList, nil
	}, nil
}
*/

/*
var personType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Person",
	Fields: graphql.Fields{
		"uri":          &graphql.Field{Type: graphql.String},
		"id":           &graphql.Field{Type: graphql.String},
		"sourceId":     &graphql.Field{Type: graphql.String},
		"primaryTitle": &graphql.Field{Type: graphql.String},
		"name":         &graphql.Field{Type: personNameType},
		"image":        &graphql.Field{Type: personImageType},
		"type":         &graphql.Field{Type: personTypeType},
		// NOTE: these are 'fake' paging - all data is still parsed in json parse
		"overviewList": &graphql.Field{Type: graphql.NewList(overviewType)},
		"keywordList":  &graphql.Field{Type: graphql.NewList(keywordType)},
		"extensionList": &graphql.Field{
			//Type: graphql.NewList(extensionType),
			Type: extensionListType,
			Args: graphql.FieldConfigArgument{
				"size": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 100},
				"from": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: extensionResolver,
		},
		// these can actually be paged, since they involve further queries
		"publicationList": &graphql.Field{
			Type: publicationListType,
			Args: graphql.FieldConfigArgument{
				"size": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 100},
				"from": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: publicationResolver,
		},
		"affiliationList": &graphql.Field{
			Type: graphql.NewList(affiliationType),
			Args: graphql.FieldConfigArgument{
				"size": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 100},
				"from": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: affiliationResolver,
		},
		"educationList": &graphql.Field{
			Type: graphql.NewList(educationType),
			Args: graphql.FieldConfigArgument{
				"size": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 100},
				"from": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: educationResolver,
		},
		"grantList": &graphql.Field{
			Type: graphql.NewList(grantType),
			Args: graphql.FieldConfigArgument{
				"size": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 100},
				"from": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
			},
			Resolve: grantResolver,
		},
	},
})

*/

/*
var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"personList":      GetPeople,
		"person":          GetPerson,
		"publicationList": GetPublications,
	},
})
*/

/*
var personListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PersonList",
	Fields: graphql.Fields{
		"results":  &graphql.Field{Type: graphql.NewList(personType)},
		"pageInfo": &graphql.Field{Type: pageInfoType},
	},
})


type ExtensionList struct {
	Results  []models.Extension `json:"data"`
	PageInfo PageInfo           `json:"pageInfo"`
}
*/

/*
var extensionListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "extensionList",
	Fields: graphql.Fields{
		"results":  &graphql.Field{Type: graphql.NewList(extensionType)},
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
*/
/*
var GetPerson = &graphql.Field{
	Type:        personType,
	Description: "Get Person",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		var person = models.Person{}
		ctx := context.Background()
		client = GetClient()

		id := params.Args["id"].(string)
		log.Printf("looking for person %s\n", id)

		get1, err := client.Get().
			Index("people").
			Id(id).
			Do(ctx)
		if err != nil {
			return person, err
		}

		// FIXME: in order to page something (like extensionList) that is part of 'Person'
		// would need to change this - otherwise it's fake paging
		err = json.Unmarshal(*get1.Source, &person)

		if err != nil {
			return person, err
		}
		return person, nil
	},
}

*/

/*
type ExtensionList struct {
	Results  []models.Extension `json:"data"`
	PageInfo PageInfo           `json:"pageInfo"`
}

type PageInfo struct {
	PerPage    int `json:"perPage"`
	Page       int `json:"page"`
	TotalPages int `json":totalPages"`
}

type PersonList struct {
	Results  []models.Person `json:"data"`
	PageInfo PageInfo        `json:"pageInfo"`
}

type PublicationList struct {
	Results  []models.Publication `json:"data"`
	PageInfo PageInfo             `json:"pageInfo"`
}
*/

/*
var GetPeople = &graphql.Field{
	Type: personListType,
	//Type:        graphql.NewList(personType),
	Description: "Get all people",
	Args: graphql.FieldConfigArgument{
		"size": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 100},
		"from": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		var people []models.Person
		ctx := context.Background()
		client = GetClient()

		size := params.Args["size"].(int)
		from := params.Args["from"].(int)

		q := elastic.NewMatchAllQuery()

		searchResult, err := client.Search().
			Index("people").
			//Type().
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
			person := models.Person{}
			err := json.Unmarshal(*hit.Source, &person)
			if err != nil {
				panic(err)
			}
			people = append(people, person)
		}

		pageInfo := PageInfo{PerPage: size,
			Page:       (from / size) + 1,
			TotalPages: (int(searchResult.TotalHits()) / size) + 1}

		personList := PersonList{Results: people, PageInfo: pageInfo}
		return personList, nil
	},
}

var GetPublications = &graphql.Field{
	Type:        graphql.NewList(publicationType),
	Description: "Get all publications",
	Args: graphql.FieldConfigArgument{
		"size": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 100},
		"from": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 1},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		var publications []models.Publication
		ctx := context.Background()
		// should query elastic here
		client = GetClient()

		size := params.Args["size"].(int)
		from := params.Args["from"].(int)

		q := elastic.NewMatchAllQuery()

		searchResult, err := client.Search().
			Index("publications").
			//Type().
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
			publication := models.Publication{}
			err := json.Unmarshal(*hit.Source, &publication)
			if err != nil {
				panic(err)
			}
			publications = append(publications, publication)
		}
		return publications, nil
	},
}
*/
type Config struct {
	Elastic elasticSearch `toml:"elastic"`
	Graphql graphqlServer `toml:"graphql"`
}

type elasticSearch struct {
	Url string
}

type graphqlServer struct {
	Port int
}

var conf Config

/*
var client *elastic.Client

func GetClient() *elastic.Client {
	return client
}
*/

func main() {
	var err error
	var configFile string

	log.SetOutput(os.Stdout)

	flag.StringVar(&configFile, "config", "./config.toml", "a config filename")

	flag.Parse()

	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		fmt.Println("could not find config file, use -c option")
		os.Exit(1)
	}

	// NOTE: elastic client is supposed to be long-lived
	// see https://github.com/olivere/elastic/blob/release-branch.v6/client.go
	//client, err = elastic.NewClient(elastic.SetURL(conf.Elastic.Url), elastic.SetSniff(false))
	graphql_endpoint.Client, err = elastic.NewClient(elastic.SetURL(conf.Elastic.Url), elastic.SetSniff(false))

	if err != nil {
		panic(err)
	}

	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql_endpoint.RootQuery,
	})

	c := cors.New(cors.Options{
		AllowCredentials: true,
	})

	h := handler.New(&handler.Config{
		Schema:   &schema,
		GraphiQL: true,
		Pretty:   true,
	})

	http.Handle("/graphql", c.Handler(h))

	// NOTE: if not configured this would default to 0
	var port = 9001
	if conf.Graphql.Port > 0 {
		port = conf.Graphql.Port
	}

	portConfig := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portConfig, nil)
}
