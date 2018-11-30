defmodule GraphqlEndpointWeb.Schema.Types do
  use Absinthe.Schema.Notation

  @desc """
  Anything within streamer has to have an access_token specified to pull back the information.
  """
  object :person do
    field(:uri, :string)
    field(:first_name, :string)
  end
end
