import { ApolloClient } from 'apollo-client'
import { HttpLink } from 'apollo-link-http'
import { InMemoryCache } from 'apollo-cache-inmemory'
import fetch from 'node-fetch'

const client = new ApolloClient({
  link: new HttpLink({
    fetch,
    uri: 'http://localhost:4001/api',
  }),
  cache: new InMemoryCache(),
})

export default client
