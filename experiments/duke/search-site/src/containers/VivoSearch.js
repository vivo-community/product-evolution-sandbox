import React, { Component } from 'react'
import { Provider } from 'react-redux'
import { Route, Switch } from 'react-router' // react-router v4
import { ConnectedRouter } from 'connected-react-router'
import VivoSearchApp from './VivoSearchApp'

//import { sagaMiddleware, configureStoreSaga } from '../configureStore'

import { history } from '../reducers/search'

import { store } from '../configureStore'

//const store = configureStoreSaga()
//import rootSaga from '../actions/sagas'

//store.runSaga = sagaMiddleware.run
//store.runSaga(rootSaga)

export default class VivoSearch extends Component {
  
  constructor(props) {
    super(props)
  }

  render() {
    return (
      <Provider store={store}>
        <ConnectedRouter history={history}>
          <Switch>
           <Route exact path='/' component={VivoSearchApp}/>
          </Switch>
        </ConnectedRouter>
      </Provider>
    )        
  }
}


