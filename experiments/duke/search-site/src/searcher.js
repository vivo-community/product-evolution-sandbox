import elasticsearch from 'elasticsearch'
import PAGE_ROWS from './actions/types'

class ElasticQuery {
 
  constructor(host) {
    var client = new elasticsearch.Client({
      host: host,
      log: 'trace'
    })

    this.host = host
    this.client = client
    this._term = ''
    this._start = 0
  }

  execute() {
    // TODO: highlight not actually returning anything right now
    // also, would probably build this with a 'querybuilder' type
    // of thing
    return this.client.search({
      index: 'people, publications, grants, educations, affiliations',
      body: {
       "from" : `${this.start}`, "size" : PAGE_ROWS,
        query: {
           "query_string": { query: `${this.term}*` }
        },
        "highlight" : {
            "fields" : {
                "_all" : {}
            }
        }
      }
    })

  }
  
  // not sure there's a good reason for this get/set indirection
  set term(term) {
    this._term = term
    return this  
  }

  get term() {
    return this._term
  }
  
  set start(start) {
    this._start = start
    return this  
  }

  get start() {
    return this._start
  }

}

export default ElasticQuery

