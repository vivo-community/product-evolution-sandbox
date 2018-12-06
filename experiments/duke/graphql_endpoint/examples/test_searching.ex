alias GraphqlEndpoint.Search

query = %{query: %{"match_all" => %{}}, size: 10}

Search.fetch("affiliations", ["affiliation"], query)
|> IO.inspect()
