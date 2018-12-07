# GraphqlEndpoint setup (new instructions)

`cd product-evolution-sandbox/experiments/duke/graphql_endpoint`
`docker-compse up graphql`

Note: if you receive the error "elastic_import_es_network declared as external, but could not be found" you may need to change the spelling of the name in the last line of docker-compose.yml to "elasticimport_es_network"

The first time you run this command, the necessary files to run the endpoint should be downloaded and installed. Once that is complete, pe_graphql should start up automatically.

GraphiQL can be viewed in the browser at http://docker:4000/api/graphiql


# GraphqlEndpoint (original instructions)

To start your Phoenix server:

  * Install dependencies with `mix deps.get`
  * Start Phoenix endpoint with `mix phx.server`

Now you can visit [`localhost:4000`](http://localhost:4000) from your browser.

Ready to run in production? Please [check our deployment guides](https://hexdocs.pm/phoenix/deployment.html).

## Learn more

  * Official website: http://www.phoenixframework.org/
  * Guides: https://hexdocs.pm/phoenix/overview.html
  * Docs: https://hexdocs.pm/phoenix
  * Mailing list: http://groups.google.com/group/phoenix-talk
  * Source: https://github.com/phoenixframework/phoenix