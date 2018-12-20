package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/OIT-ads-web/graphql_endpoint/models"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/olivere/elastic"
	"net/http"
	"os"
)

type Config struct {
	Elastic elasticSearch `toml:"elastic"`
}

type elasticSearch struct {
	Url string
}

var client *elastic.Client
var conf Config

func GetClient() *elastic.Client {
	return client
}

/*
type PersonKeyword struct {
	Uri   string `json:"uri"`
	Label string `json:"label"`
}

type PersonImage struct {
	Main      string `json:"main"`
	Thumbnail string `json:"thumbnail"`
}

type PersonName struct {
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	MiddleName *string `json:"middleName"`
}

type PersonType struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type OverviewType struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type PersonOverview struct {
	Label string       `json:"overview"`
	Type  OverviewType `json:"type"`
}

type Extension struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Person struct {
	Id           string           `json:"id"`
	Uri          string           `json:"uri"`
	SourceId     string           `json:"sourceId"`
	PrimaryTitle string           `json:"primaryTitle"`
	Name         PersonName       `json:"name" elastic:"type:object"`
	Image        PersonImage      `json:"image" elastic:"type:object"`
	Type         PersonType       `json:"type" elastic:"type:object"`
	OverviewList []PersonOverview `json:"overviewList" elastic:"type:nested"`
	KeywordList  []PersonKeyword  `json:"keywordList" elastic:"type:nested"`
	Extensions   []Extension      `json:"extensions" elastic:"type:nested"`
}
*/

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
		"id":  &graphql.Field{Type: graphql.String},
		"uri": &graphql.Field{Type: graphql.String},
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

// var personResolver = func(params graphql.ResolveParams) (interface{}, error) {
// }	
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
		"publicationList": &graphql.Field{
			Type: graphql.NewList(authorshipType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				person, _ := params.Source.(models.Person)
				//var authorships []models.Authorship
				var publications []models.Publication

				ctx := context.Background()
				client = GetClient()

				q := elastic.NewMatchQuery("personId", person.Id)

				searchResult, err := client.Search().
					Index("authorships").
					Query(q).
					From(0).
					Size(1000).
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
		            get1, err := client.Get().
			            Index("publications").
			            Id(publicationId).
			            Do(ctx)
					
					publication := models.Publication{}
					err = json.Unmarshal(*get1.Source, &publication)

					if err != nil {
						panic(err)
					}
					publications = append(publications, publication)
				}
				return publications, nil
			},
		},
		/*
		"authorships": &graphql.Field{
			Type: graphql.NewList(authorshipType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				person, _ := params.Source.(models.Person)
				var authorships []models.Authorship
				ctx := context.Background()
				client = GetClient()

				q := elastic.NewMatchQuery("personId", person.Id)

				searchResult, err := client.Search().
					Index("authorships").
					Query(q).
					From(0).
					Size(1000).
					Do(ctx)
				if err != nil {
					// Handle error
					panic(err)
				}

				//TotalHits()

				// load publications - with authorships?
				for _, hit := range searchResult.Hits.Hits {
					authorship := models.Authorship{}
					err := json.Unmarshal(*hit.Source, &authorship)
					if err != nil {
						panic(err)
					}
					authorships = append(authorships, authorship)

				}
				return authorships, nil
			},
		},
		*/
	},
})

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"personList":      GetPeople,
		"person":          GetPerson,
		"publicationList": GetPublications,
	},
})

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
	},
}

var GetPeople = &graphql.Field{
	Type:        graphql.NewList(personType),
	Description: "Get all people",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		var people []models.Person
		ctx := context.Background()
		// should query elastic here
		client = GetClient()

		q := elastic.NewMatchAllQuery()

		searchResult, err := client.Search().
			Index("people").
			//Type().
			Query(q).
			From(0).
			Size(1000).
			//Pretty(true).
			// Timeout("1000ms"). or
			// Timeout(1000).
			Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}

		//TotalHits()

		for _, hit := range searchResult.Hits.Hits {
			person := models.Person{}
			err := json.Unmarshal(*hit.Source, &person)
			if err != nil {
				panic(err)
			}
			people = append(people, person)
		}
		return people, nil
	},
}

var GetPublications = &graphql.Field{
	Type:        graphql.NewList(publicationType),
	Description: "Get all publications",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		var publications []models.Publication
		ctx := context.Background()
		// should query elastic here
		client = GetClient()

		q := elastic.NewMatchAllQuery()

		searchResult, err := client.Search().
			Index("publications").
			//Type().
			Query(q).
			From(0).
			Size(1000).
			//Pretty(true).
			// Timeout("1000ms"). or
			// Timeout(1000).
			Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}

		//TotalHits()

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

func main() {
	var err error
	var configFile string
	flag.StringVar(&configFile, "config", "./config.toml", "a config filename")

	flag.Parse()

	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		fmt.Println("could not find config file, use -c option")
		os.Exit(1)
	}

	// NOTE: elastic client is supposed to be long-lived
	// see https://github.com/olivere/elastic/blob/release-branch.v6/client.go
	client, err = elastic.NewClient(elastic.SetURL(conf.Elastic.Url))
	if err != nil {
		panic(err)
	}

	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: RootQuery,
		//Mutation: RootMutation,
	})

	h := handler.New(&handler.Config{
		Schema:   &schema,
		GraphiQL: true,
		Pretty:   true,
	})

	http.Handle("/graphql", h)
	http.ListenAndServe(":9001", nil)
}
