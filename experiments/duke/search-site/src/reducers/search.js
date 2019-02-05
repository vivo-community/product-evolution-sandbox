//http://spapas.github.io/2016/03/02/react-redux-tutorial/#components-notification-js
import * as types from '../actions/types'

// could call it #search, just called it #searchReducer to be explicit about the key name
// in the combineReducers method
function searchReducer(search = { isFetching: false, results: {}}, action) {
  console.debug("searchReducer called " + action)
  switch (action.type) {
  case types.REQUEST_SEARCH:
    console.debug("REQUEST_SEARCH")
    console.debug(action.searchFields)
   
    return { ...search, 
      isFetching: true,
      results: action.results,
      searchFields: action.searchFields,
      lastUpdated: action.requestedAt
  }
  case types.RECEIVE_SEARCH:
    console.debug("RECEIVE_SEARCH")
    console.debug(action.results)
    return { ...search, 
      isFetching: false,
      results: action.results,
      lastUpdated: action.receivedAt
  }

  case types.SEARCH_FAILED:
    console.debug("SEARCH_FAILED")   
    return { ...search,
      isFetching: false,
      message: action.message,
      lastUpdated: action.failedAt
  }
  default:
    console.debug("default")
    return search;
  }
}


import { connectRouter } from 'connected-react-router'
import { combineReducers } from 'redux'
//import { routerReducer  } from 'react-router-redux'
import { createBrowserHistory } from 'history'
export const history = createBrowserHistory()

// NOTE: each reducer combines to effect the global state,
// but only the named one - so, in effect, it's like
// a set of named sub-states within the global state
// e.g. state = {'search': .., 'routing': .. 'init': .. }
// it's not necessary to explicitly set the key like this,
// as it will default to the name of the function.  I just
// did it to be obvious
const mainReducer = combineReducers({
  search: searchReducer,
  router: connectRouter(history),
})

export default {
  mainReducer,
  searchReducer
}
