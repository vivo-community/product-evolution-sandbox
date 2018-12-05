defmodule GraphqlEndpointWeb.Resolvers.People do
  alias GraphqlEndpoint.Search

  def fetch_person(uri) do
    q = %{query: %{match_all: %{}}, size: 10}
    Search.fetch("people", ["person"], q)
  end
end
