alias GraphqlEndpointWeb.Resolvers.People

# IO.inspect(People.fetch_by_id("per1709582"))

#alias GraphqlEndpointWeb.Resolvers.Affiliations
#IO.inspect(Affiliations.fetch(%{}, %{person_id: "per4142272"} ,%{}))

#alias GraphqlEndpointWeb.Resolvers.Grants
#IO.inspect(Grants.fetch(%{}, %{person_id: "per7738882"} ,%{}))

alias GraphqlEndpointWeb.Resolvers.Publications
IO.inspect(Publications.fetch(%{}, %{person_id: "per0120862"} ,%{}))
