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

defmodule GraphqlEndpointWeb.Schema.Types do
  use Absinthe.Schema.Notation
  alias GraphqlEndpointWeb.Resolvers

  @desc """
  A person
  """
  object :person do
    field(:id, :string)
    field(:uri, :string)
    field(:image, :image)
    field(:name, :name)
    field(:overview_list, list_of(:overview))

    field(:affiliation_list, list_of(:affiliation)) do
      resolve(&Resolvers.Affiliations.fetch/3)
    end

    field(:education_list, list_of(:education)) do
      resolve(&Resolvers.Educations.fetch/3)
    end

    field(:grant_list, list_of(:grant)) do
      resolve(&Resolvers.Grants.fetch/3)
    end

    field(:publication_list, list_of(:publication)) do
      resolve(&Resolvers.Publications.fetch/3)
    end
  end

  object :image do
    field(:main, :string)
    field(:thumbnail, :string)
  end

  object :name do
    field(:first_name, :string)
    field(:last_name, :string)
    field(:middle_name, :string)
  end

  object :overview do
    field(:overview, :string)
    field(:type, :type)
  end

  object :type do
    field(:code, :string)
    field(:label, :string)
  end

  object :affiliation do
    field(:id, :string)
    field(:label, :string)
    field(:start_date, :date_resolution)
  end

  object :education do
    field(:label, :string)
    field(:org, :organization)
  end

  object :organization do
    field(:id, :string)
    field(:label, :string)
  end

  object :date_resolution do
    field(:date_time, :string)
    field(:resolution, :string)
  end

  # object :funding_role do
  # field(:date_time, :string)
  # field(:label, :string)
  # end

  # object :authorship do
  # field(:date_time, :string)
  # field(:resolution, :string)
  # end

  object :grant do
    field(:id, :string)
    field(:label, :string)
    field(:role_name, :string)
    field(:start_date, :date_resolution)
    field(:end_date, :date_resolution)
  end

  object :venue do
    field(:uri, :string)
    field(:label, :string)
  end

  object :publication do
    field(:id, :string)
    field(:author_list, :string)
    field(:doi, :string)
    field(:label, :string)
    field(:role_name, :string)
    field(:venue, :venue)
  end

  object :publication_list do
    field(:results, list_of(:publication))
    field(:page_info, :page_info)
  end

  object :page_info do
    field(:per_page, :integer)
    field(:page, :integer)
    field(:total_pages, :integer)
  end
end

*/

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

var grantType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Grant",
	Fields: graphql.Fields{
		"id":    &graphql.Field{Type: graphql.String},
		"label": &graphql.Field{Type: graphql.String},
		"roleName": &graphql.Field{Type: graphql.String},
		"startDate": &graphql.Field{Type: dateResolutionType},
		"endDate": &graphql.Field{Type: dateResolutionType},
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
		"overviewList": &graphql.Field{Type: graphql.NewList(overviewType)},
		"publicationList": &graphql.Field{
			Type: graphql.NewList(publicationType),
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
		"affiliationList": &graphql.Field{
			Type: graphql.NewList(affiliationType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				person, _ := params.Source.(models.Person)
				var affiliations []models.Affiliation
				ctx := context.Background()
				client = GetClient()

				q := elastic.NewMatchQuery("personId", person.Id)

				searchResult, err := client.Search().
					Index("affiliations").
					Query(q).
					From(0).
					Size(1000).
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
				return affiliations, nil
			},
		},
		"educationList": &graphql.Field{
			Type: graphql.NewList(educationType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				person, _ := params.Source.(models.Person)
				var educations []models.Education
				ctx := context.Background()
				client = GetClient()

				q := elastic.NewMatchQuery("personId", person.Id)

				searchResult, err := client.Search().
					Index("educations").
					Query(q).
					From(0).
					Size(1000).
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
				return educations, nil
			},
		},
		"grantList": &graphql.Field{
			Type: graphql.NewList(grantType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				person, _ := params.Source.(models.Person)
				var grants []models.Grant

				ctx := context.Background()
				client = GetClient()

				q := elastic.NewMatchQuery("personId", person.Id)

				searchResult, err := client.Search().
					Index("funding-roles").
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
					fundingRole := models.FundingRole{}
					err := json.Unmarshal(*hit.Source, &fundingRole)
					if err != nil {
						panic(err)
					}

					grantId := fundingRole.GrantId
					get1, err := client.Get().
						Index("grants").
						Id(grantId).
						Do(ctx)

					grant := models.Grant{}
					err = json.Unmarshal(*get1.Source, &grant)

					if err != nil {
						panic(err)
					}
					grants = append(grants, grant)
				}
				return grants, nil
			},
		},
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
