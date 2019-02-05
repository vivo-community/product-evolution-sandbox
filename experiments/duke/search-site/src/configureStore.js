import { createStore, applyMiddleware } from 'redux'
import { createLogger } from 'redux-logger'
import createSagaMiddleware from 'redux-saga'
const sagaMiddleware = createSagaMiddleware()
const loggerMiddleware = createLogger()

import { routerMiddleware } from 'connected-react-router'
import reducers from './reducers/search'

let middlewares = [sagaMiddleware, routerMiddleware]

if (process.env.NODE_ENV != 'production') {
  middlewares.push(loggerMiddleware)
}

// FIXME: don't like all these versions of basically the same thing
const createStoreWithMiddleware = applyMiddleware(
  ...middlewares
)(createStore)

const initialState = {}

const store = createStore(reducers.mainReducer,
  applyMiddleware(sagaMiddleware)
)

import rootSaga from './actions/sagas'
sagaMiddleware.run(rootSaga)

//function configureStoreSaga(initialState = initialState) {
//  return createStoreWithMiddleware(reducers.mainReducer, initialState)
//}
//const store = configureStoreSaga()
//import rootSaga from '../actions/sagas'

//store.runSaga = sagaMiddleware.run
//store.runSaga(rootSaga)


export { store }

