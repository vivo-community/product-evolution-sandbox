import { client } from './config'
import { inspect, error } from './search'

client.search({
  index: 'people',
  //type: 'person',
  body: {
    query: {
      match: {
        'name.lastName': 'Swaminathan'
        //'keywordList.label': 'Cardiomyopathies'
        //'_all': 'Swaminathan'
      }
    }
  }
}).then(inspect, error);



