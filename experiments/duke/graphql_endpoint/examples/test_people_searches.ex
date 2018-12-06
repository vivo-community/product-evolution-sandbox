alias GraphqlEndpointWeb.Resolvers.People

# IO.inspect(People.fetch_by_id("per1709582"))

alias GraphqlEndpointWeb.Resolvers.Affiliations

IO.inspect(Affiliations.fetch(%{}, %{person_id: "per4142272"} ,%{}))