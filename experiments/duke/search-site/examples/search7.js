import { client } from './config'
import { inspectAgg, error } from './search'

//https://stackoverflow.com/questions/48784733/aggregation-on-keyword-field-fails

/*
       "experience_slug": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        },
        "analyzer": "default",
        "search_analyzer": "default_search"
      },
      */
// facets=aggregation, bucketing etc...
client.search({
  //size: "0",
  index: ["people"],
  //index: "people",
  //type: 'person',
  body: {
    query: {
      "query_string": { query: "model*" }
    },
    aggs: {
       "typeFilter": {
           //"terms" : { "field": "keywordList.label" } 
   
       /*
       "typeFilter": {
           "terms" : { "field": "type" } 
       */
         
         "nested": {
             "path": "keywordList"
         },
         "aggs": {
           "keywords_filter" : { "terms" : { "field": "keywordList.label.keyword" } }
         }
         
      },
      /*
      "keywordFilter": {
           "terms" : { "field": "keywordList" } 
      }
      */
    }
  }
}).then(inspectAgg, error);



