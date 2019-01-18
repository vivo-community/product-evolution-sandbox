import { client } from './config'
import { inspect, error } from './search'

client.search({
  //index: ["people", "publications"],
  //type: 'person',
  //q: '*:Nelson'
  body: {
    "query": {
      "multi_match": {
        "query": "Swaminathan",
        "type": "cross_fields",
        //"fields": "name.firstName"
      }
    }
  }
}).then(inspect, error);



