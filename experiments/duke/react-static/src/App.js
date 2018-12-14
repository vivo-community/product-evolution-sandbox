import React from 'react'
import { Root, Routes } from 'react-static'
import { Link } from '@reach/router'
import { ApolloProvider } from 'react-apollo'
import client from './connectors/apollo'
import './app.css'

const App = () => (
  <ApolloProvider client={client}>
    <Root>
      <div className="navigation">
        <nav>
          <Link exact to="/">Home</Link>
          <Link to="/people">People</Link>
        </nav>
      </div>
      <div className="content">
        <Routes />
      </div>
    </Root>
  </ApolloProvider>
)

export default App
