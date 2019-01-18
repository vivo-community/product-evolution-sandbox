import * as types from './types'

export function requestSearch(searchFields) {
  console.debug("action=requestSearch")
  return {
    type: types.REQUEST_SEARCH,
    results: {results: {}},
    isFetching: true,
    searchFields: searchFields,
    requestedAt: Date.now()
  }
}

export function receiveSearch(data) {
  console.debug("action=receiveSearch")
  console.debug(data)
  return {
    type: types.RECEIVE_SEARCH,
    results: data,
    isFetching: false,
    receivedAt: Date.now()
  }
}

export function cancelSearch() {
  return {
    type: types.SEARCH_CANCELLED
  }
}

export function searchFailed(message) {
  return {
    type: types.SEARCH_FAILED,
    message: message,
    failedAt: Date.now()
  }
}

// allow all to be exported at once into an 'actions' object
export default {
  requestSearch,
  receiveSearch,
  cancelSearch,
  searchFailed,
} 
