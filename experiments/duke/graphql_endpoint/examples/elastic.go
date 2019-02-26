package examples

import (
	"context"
	"encoding/json"
	el "github.com/OIT-ads-web/graphql_endpoint/elastic"
	"github.com/davecgh/go-spew/spew"
	"github.com/olivere/elastic"
	"log"
)

/*
type PersonFinder struct {
	// facets
	department     []string
	type           []string
	keyword        []string
	from       int
	size       int
	sort       []string
}

type PersonFinderResponse struct {
	Total          int64
	Films          []*Film
	Genres         map[string]int64
	YearsAndGenres map[int][]NameCount // {1994: [{"Crime":1}, {"Drama":2}], ...}
}

// Example: "name" or "-year".
func (f *Finder) Sort(sort ...string) *Finder {
	if f.sort == nil {
		f.sort = make([]string, 0)
	}
	f.sort = append(f.sort, sort...)
	return f
}

// sorting applies sorting to the service.
func (f *Finder) sorting(service *elastic.SearchService) *elastic.SearchService {
	if len(f.sort) == 0 {
		// Sort by score by default
		service = service.Sort("_score", false)
		return service
	}

	// Sort by fields; prefix of "-" means: descending sort order.
	for _, s := range f.sort {
		s = strings.TrimSpace(s)

		var field string
		var asc bool

		if strings.HasPrefix(s, "-") {
			field = s[1:]
			asc = false
		} else {
			field = s
			asc = true
		}

		// Maybe check for permitted fields to sort

		service = service.Sort(field, asc)
	}
	return service
}


func (f *Finder) aggs(service *elastic.SearchService) *elastic.SearchService {
	// Terms aggregation by genre
	agg := elastic.NewTermsAggregation().Field("type.label")
	service = service.Aggregation("types", agg)

	nested := NewNestedAggregation().Path("keywordList")
	//"keyword" : { "terms" : { "field": "keywordList.label.keyword" } }
	nested = nested.SubAggregations("keyword", NewTermsAggregation().Field("keywordList.label.keyword")


	// Add a terms aggregation of Year, and add a sub-aggregation for Genre
	subAgg := elastic.NewTermsAggregation().Field("keywords")
	service = service.Aggregation("keywords", nested)


	elastic.NewNestedAggregation?
	agg = elastic.NewTermsAggregation().Field("year").
		SubAggregation("genres_by_year", subAgg)
	service = service.Aggregation("years_and_genres", agg)

	return service
}

--
    aggs: {
      "types" : { "terms": {"field" : "type.label" }},
      "keywords": {
         "nested": {
             "path": "keywordList"
         },
         "aggs": {
           "keyword" : { "terms" : { "field": "keywordList.label.keyword" } }
         }
      },
      "affliliations": {
         "nested": {
             "path": "affiliationList"
         },
         "aggs": {
           //"title": { "terms" : { "field": "affiliationList.label.title" } },
           "department": { "terms" : { "field": "affiliationList.organization.label.dept" } }
         }
      }
    }

*/

// NewTermsLookup.Path("keywordList")
// TODO: will need to take filter and build search
// sort, facets etc...
// see https://github.com/olivere/elastic/wiki/QueryDSL
//
//  search := client.Search().Index(indexName).Type("_doc").Pretty(f.pretty)
//	search = f.query(search)
//	search = f.aggs(search)
//	search = f.sorting(search)
//	search = f.paginate(search)
/*

"aggregations": {
      "types": {
        "doc_count_error_upper_bound": 0,
        "sum_other_doc_count": 0,
        "buckets": [
          {
            "key": "http://vivoweb.org/ontology/core#FacultyMember",
            "doc_count": 658
          },
          {
            "key": "http://vivoweb.org/ontology/core#Student",
            "doc_count": 351
          },
          {
            "key": "http://vivoweb.org/ontology/core#NonFacultyAcademic",
            "doc_count": 117
          },
          {
            "key": "http://vivo.duke.edu/vivo/ontology/duke-extension#Affiliate",
            "doc_count": 16
          }
        ]
      },
      "keywords": {
        "doc_count": 2436,
        "keyword": {
          "doc_count_error_upper_bound": 0,
          "sum_other_doc_count": 2349,
          "buckets": [
            {
              "key": "Humans",
              "doc_count": 13
            },
            {
              "key": "Animals",
              "doc_count": 10
            },
            {
              "key": "Computer Simulation",
              "doc_count": 10
            },
            {
              "key": "Female",
              "doc_count": 9
            },
            {
              "key": "Male",
              "doc_count": 9
            },
            {
              "key": "Adult",
              "doc_count": 8
            },
            {
              "key": "Aged",
              "doc_count": 7
            },




	      "affliliations": {
         "nested": {
             "path": "affiliationList"
         },
         "aggs": {
           //"title": { "terms" : { "field": "affiliationList.label.title" } },
           "department": { "terms" : { "field": "affiliationList.organization.label.dept" } }
         }
      }
    }

*/

func ExampleIdQuery() {
	el.IdQuery("people", []string{"per4774112", "per8608642"})
}

func ExampleListPeople() {
	el.ListAll("people")
}

func ExampleAggregations() {
	ctx := context.Background()
	client := el.GetClient()

	q := elastic.NewMatchAllQuery()

	service := client.Search().
		Index("people").
		Query(q).
		From(0).
		Size(10)

	agg := elastic.NewTermsAggregation().Field("type.label")
	service = service.Aggregation("types", agg)

	nested := elastic.NewNestedAggregation().Path("keywordList")
	subAgg := nested.SubAggregation("keyword", elastic.NewTermsAggregation().Field("keywordList.label.keyword"))

	nested2 := elastic.NewNestedAggregation().Path("affiliationList")
	subAgg2 := nested2.SubAggregation("department",
		elastic.NewTermsAggregation().Field("affiliationList.organization.label.dept"))

	service = service.Aggregation("keywords", subAgg)
	service = service.Aggregation("affiliations", subAgg2)

	searchResult, err := service.Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	for _, hit := range searchResult.Aggregations {
		var obj interface{}
		err := json.Unmarshal(*hit, &obj)
		if err != nil {
			panic(err)
		}
		str := spew.Sdump(obj)
        log.Println(str)
	}
}
