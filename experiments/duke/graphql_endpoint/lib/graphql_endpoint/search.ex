defmodule GraphqlEndpoint.Search do
  def fetch(index, types, query) do
    Elastix.Search.search(endpoint(), index, types, query)
  end

  def endpoint do
    Application.get_env(:graphql_endpoint, :endpoint)
  end
end
