import { client } from './config'
import { inspect, error } from './search'

client.search({
  index: ["people", "publications"],
  //index: "people",
  //type: 'person',
  body: {
    query: {
      has_child: {
        'type': 'keyword',
        query: {
           "query_string": { query: "Cardiomyopathies*" }
        }
      }
    }
  }
}).then(inspect, error);



