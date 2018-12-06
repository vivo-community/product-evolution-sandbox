defmodule GraphqlEndpointWeb.Resolvers.Educations do
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

    case Search.fetch("educations", ["education"], q) do
      {:ok, %HTTPoison.Response{status_code: 200, body: body}} ->
        {:ok, process_body(body)}

      {:error, %HTTPoison.Error{reason: reason}} ->
        {:error, reason}
    end
  end

  defp process_body(%{"hits" => _h = %{"hits" => hits}}) do
    Enum.map(hits, &process_education/1)
  end

  defp process_body(body) do
    IO.inspect(body, label: "FULL BODY")
    %{"error" => "unable to process request"}
  end

  defp process_education(es_object) do
    es_object
    |> Map.get("_source")
    |> IO.inspect(label: "object")
    |> JsonHelper.atomize_understore_keys()
  end
end
