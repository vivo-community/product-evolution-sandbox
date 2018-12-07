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
  end
end
