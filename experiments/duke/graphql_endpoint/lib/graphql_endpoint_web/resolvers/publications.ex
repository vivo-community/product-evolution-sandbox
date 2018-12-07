defmodule GraphqlEndpointWeb.Resolvers.Publications do
  alias GraphqlEndpoint.Search
  alias GraphqlEndpoint.JsonHelper

  def fetch(_parent, %{person_id: person_id}, _context) do
    fetch_by_person_id(person_id)
  end

  def fetch(parent=%{id: person_id}, _args, _context) do
    IO.puts(">>>> inside of fetch: ")
    parent
    |> IO.inspect(label: "parent")
    fetch_by_person_id(person_id)
  end

  def fetch(_, _, _) do
    {:error, "Unable to retrieve"}
  end

  defp fetch_by_person_id(id) do
    q = %{query: %{match: %{"personId" => id}}, size: 10}

    case Search.fetch("authorships", ["authorship"], q) do
      {:ok, %HTTPoison.Response{status_code: 200, body: body}} ->
        {:ok, process_body(body)}

      {:error, %HTTPoison.Error{reason: reason}} ->
        {:error, reason}
    end
  end

  defp process_body(%{"hits" => _h = %{"hits" => hits}}) do
    Enum.map(hits, &process_authorship/1)
  end

  defp process_body(body) do
    IO.inspect(body, label: "FULL BODY")
    %{"error" => "unable to process request"}
  end

  defp process_authorship(es_object) do
    # needs to get underlying publication
    authorship = Map.get(es_object, "_source")
    publication_id = Map.get(authorship, "publicationId")
    q = %{query: %{match: %{"id" => publication_id}}, size: 100}

    publication = case Search.fetch("publications", ["publication"], q) do
      {:ok, %HTTPoison.Response{status_code: 200, body: body}} ->
        {:ok, process_publication_body(body)}
      {:error, %HTTPoison.Error{reason: reason}} ->
        {:error, reason}
    end

    # probably don't need to get this twice
    authorship = authorship
    |> JsonHelper.atomize_understore_keys()

    # just assuming one
    publication_match = Enum.at(elem(publication, 1), 0)

    compound = %{
      id: publication_match[:id],
      label: publication_match[:label],
      author_list: publication_match[:author_list],
      doi: publication_match[:doi],
      venue: publication_match[:venue],
      role_name: authorship[:label]
    }
    compound
  end

  defp process_publication_body(%{"hits" => _h = %{"hits" => hits}}) do
    Enum.map(hits, &process_publication/1)
  end

  defp process_publication_body(body) do
    IO.inspect(body, label: "FULL BODY")
    %{"error" => "unable to process request"}
  end

  defp process_publication(es_object) do
    es_object
    |> Map.get("_source")
    |> JsonHelper.atomize_understore_keys()
  end

end
