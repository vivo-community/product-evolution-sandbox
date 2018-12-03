defmodule GraphqlEndpointWeb.Schema do
  use Absinthe.Schema

  import_types(Absinthe.Type.Custom)
  import_types(GraphqlEndpointWeb.Schema.Types)
  # import_types BlogWeb.Schema.ContentTypes

  query do
    @desc """
    Retrieve a person based on their URI.
    """
    field :person, :person do
      arg(:uri, non_null(:string), description: "URI of the person")
    end
  end
end
