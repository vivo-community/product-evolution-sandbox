import React from 'react'
import { Router, Link } from 'react-static'
import { hot } from 'react-hot-loader'
//
import Routes from 'react-static-routes'
//
import { ApolloProvider } from 'react-apollo'
import client from './connectors/apollo'

import './app.css'

const App = () => (
  <ApolloProvider client={client}>
    <Router>
      <div>
        <nav>
          <Link exact to="/">Home</Link>
          <Link to="/people">People</Link>
          <Link to="/about">About</Link>
        </nav>
        <div className="content">
          <Routes />
        </div>
      </div>
    </Router>
  </ApolloProvider>
)

export default hot(module)(App)
