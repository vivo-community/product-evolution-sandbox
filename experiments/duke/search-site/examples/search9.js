import { client } from './config'
import { inspectAgg, error } from './search'

//https://stackoverflow.com/questions/48784733/aggregation-on-keyword-field-fails

// facets=aggregation, bucketing etc...
client.search({
  index: ["publications"],
  body: {
    query: {
      "query_string": { query: "*" }
    },
    aggs: {
      "types" : { "terms": {"field" : "type.label" }},
      "keywords": {
         "nested": {
             "path": "keywordList"
         },
         "aggs": {
           "keyword" : { "terms" : { "field": "keywordList.label.keyword" } }
         }
      }
    }
  }
}).then(inspectAgg, error);



