defmodule GraphqlEndpointWeb.Schema do
  use Absinthe.Schema

  import_types(Absinthe.Type.Custom)
  import_types(GraphqlEndpointWeb.Schema.Types)

  alias GraphqlEndpointWeb.Resolvers

  query do
    @desc """
    Retrieve a person based on their ID.
    """
    field :person, :person do
      arg(:id, non_null(:string), description: "ID of the person")
      resolve(&Resolvers.People.fetch/3)
    end

    @desc """
    Retrieve list of all people
    """
    field :person_list, list_of(:person) do
      resolve(&Resolvers.People.all/2)
    end

    @desc "Retrieve a list of publications"
    field :publication_list, :publication_list do
      arg(:size, :integer, description: "Number of publications to return")
      arg(:from, :integer, description: "Starting point for the publications (like offset)")
      arg(:query, :string, description: "Search for term in publications")
      resolve(&Resolvers.Publications.all/2)
    end
  end
end
