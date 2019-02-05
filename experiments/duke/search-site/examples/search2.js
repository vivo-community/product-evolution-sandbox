import { client } from './config'
import { inspect, error } from './search'

client.search({
  index: 'people,publications,grants',
  //type: 'person',
  //q: '*:Nelson'
  body: {
    "query": {
      "bool": {
        "must": [{"match_all": {}}],
        //"filter": [{"term": {"_all": "Problem"}}]
      }
    }
  }
}).then(inspect, error);



