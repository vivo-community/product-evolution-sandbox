import { client } from './config'
import { inspect, error } from './search'

client.search({
  index: ["people", "publications"],
  //index: "people",
  //type: 'person',
  body: {
    query: {
      nested: {
        'path': 'keywordList',
        query: {
           "query_string": { query: "Cardiomyopathies*" }
        }
      }
    }
  }
}).then(inspect, error);



