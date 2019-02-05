import * as types from './types'

import { call, put, fork, take, cancel, cancelled, all } from 'redux-saga/effects'
//import querystring from 'querystring'
import { receiveSearch, searchFailed } from './search'
import ElasticQuery from '../searcher'

export function fetchSearchApi(searchFields, maxRows=50) {
  //const solrUrl = process.env.SOLR_URL
  console.debug("in search api")
  console.debug(searchFields)
  let start = searchFields ? Math.floor(searchFields['start'] || 0) : 0
  let term = searchFields ? searchFields['term'] : ''

  console.debug(term)
  // FIXME: rows should probably be a parameter too 
  // (but within reason e.g. maybe a list of options [50, 100, 200] ...)
  // from, size ?
  //searcher.setupDefaultSearch(maxRows, start)
  const elasticUrl = process.env.ELASTIC_URL
  let searcher = new ElasticQuery(elasticUrl)
  console.debug("make elastic query")

  searcher.term = term
  searcher.start = start
  return searcher.execute()
}

// FIXME: how to cancel? 
// 
// cancel might look like this:
// https://yelouafi.github.io/redux-saga/docs/advanced/TaskCancellation.html

// 2. what watcher will do (see next)
export function* fetchSearch(action) {
  const { searchFields, dispatch } = action
  console.debug(action)

  try { 
    console.debug(searchFields)
    const results = yield call(fetchSearchApi, searchFields)

    console.debug("got some results")
    console.debug(results)
    //yield put.resolve(receiveSearch(results))
    yield put(receiveSearch(results))
    //yield put({ type: types.RECEIVE_SEARCH, results })
    //dispatch(receiveSearch(results))

  } catch(e) {
    console.debug("error?"+e)
    yield put(searchFailed(e.message))
  } finally {
    if (yield cancelled()) {
    }
  }  

}

// 3. watcher
function* watchForSearch() {
  while(true) {
    const action = yield take(types.REQUEST_SEARCH)
    const searchTask = yield fork(fetchSearch, action)
  }
}

// all of them, wrapped up for middleware
export default function* root() {
  yield all([
    fork(watchForSearch)
  ])
  //yield [
  //  fork(watchForSearch)
  //]
}
