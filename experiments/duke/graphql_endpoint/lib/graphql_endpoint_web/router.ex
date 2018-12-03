defmodule GraphqlEndpointWeb.Router do
  use GraphqlEndpointWeb, :router

  pipeline :api do
    plug :accepts, ["json"]
  end

  scope "/api" do
    pipe_through :api

    # forward "/graphiql", Absinthe.Plug.GraphiQL, schema: GraphqlEndpointWeb.Schema

    forward "/graphiql",
            Absinthe.Plug.GraphiQL,
            schema: GraphqlEndpointWeb.Schema,
            json_codec: Jason,
            interface: :simple

    forward "/", Absinthe.Plug, schema: GraphqlEndpointWeb.Schema, json_codec: Jason
  end
end
