import { client } from './config'
import { inspectAgg, error } from './search'

//https://stackoverflow.com/questions/48784733/aggregation-on-keyword-field-fails

// facets=aggregation, bucketing etc...
client.search({
  index: ["grants"],
  body: {
    query: {
      "query_string": { query: "*" }
    },
    aggs: {
      "types" : { "terms": {"field" : "type.label" }},
      "investigators": {
         "nested": {
             "path": "investigatorList"
         },
         "aggs": {
           "investigator" : { "terms" : { "field": "investigatorList.label.name" } }
         }
      }
    }
  }
}).then(inspectAgg, error);



