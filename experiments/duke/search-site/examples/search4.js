import { client } from './config'
import { inspect, error } from './search'

client.search({
  index: ["people", "publications", "grants", "affiliations", "educations"],
  //type: 'person',
  //q: '*:Nelson'
  body: {
    "query": { 
      //"query_string": { query: "Swami*" },
      "query_string": { query: "Cardiomyopathies*" }
      //"analyzeWildcard": true
    }
  }
}).then(inspect, error);



