package examples

import (
	"context"
	"encoding/json"
	"log"

	el "github.com/OIT-ads-web/graphql_endpoint/elastic"
	"github.com/davecgh/go-spew/spew"
	"github.com/olivere/elastic"
)

// should eventually try stuff more like here:
// https://github.com/olivere/elastic/wiki/QueryDSL
func ExampleIdQuery() {
	el.IdQuery("people", []string{"per4774112", "per8608642"})
}

func ExampleListPeople() {
	el.ListAll("people")
}

/*
target:

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
}
*/
/*
getting this (kind of same):

{
   "aggregations":{
      "keywords":{
         "aggregations":{
            "keyword":{
               "terms":{
                  "field":"keywordList.label.keyword"
               }
            }
         },
         "nested":{
            "path":"keywordList"
         }
      },
      "types":{
         "terms":{
            "field":"type.label"
         }
      }
   },
   "query":{
      "match_all":{

      }
   }
}
*/
func ExampleAggregations() {
	ctx := context.Background()
	client := el.GetClient()

	q := elastic.NewMatchAllQuery()

	service := client.Search().
		Index("people").
		Query(q)

	agg := elastic.NewTermsAggregation().Field("type.label")
	service = service.Aggregation("types", agg)

	nested := elastic.NewNestedAggregation().Path("keywordList")
	subAgg := nested.SubAggregation("keyword",
		elastic.NewTermsAggregation().Field("keywordList.label.keyword"))

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

	log.Println("************")
}
