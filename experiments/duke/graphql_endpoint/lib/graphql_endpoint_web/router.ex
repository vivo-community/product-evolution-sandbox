defmodule GraphqlEndpointWeb.Router do
  use GraphqlEndpointWeb, :router

  pipeline :api do
    plug :accepts, ["json"]
  end

  scope "/api", GraphqlEndpointWeb do
    pipe_through :api

    forward "/graphiql", Absinthe.Plug.GraphiQL, schema: BlogWeb.Schema

    forward "/", Absinthe.Plug, schema: BlogWeb.Schema
  end
end
