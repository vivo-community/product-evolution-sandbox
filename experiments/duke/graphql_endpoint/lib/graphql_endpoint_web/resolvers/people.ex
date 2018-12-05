defmodule GraphqlEndpointWeb.Resolvers.People do
  alias GraphqlEndpoint.Search
  alias GraphqlEndpoint.JsonHelper

  def fetch(_parent, %{id: id}, _context) do
    fetch_by_id(id)
  end

  def fetch(_, _, _) do
    {:error, "Unable to retrieve"}
  end

  def fetch_by_id(id) do
    q = %{query: %{match: %{id: id}}, size: 10}

    case Search.fetch("people", ["person"], q) do
      {:ok, %HTTPoison.Response{status_code: 200, body: body}} ->
        {:ok, process_body(body)}

      {:error, %HTTPoison.Error{reason: reason}} ->
        {:error, reason}
    end
  end

  def process_body(%{"hits" => _h = %{"hits" => hits, "total" => 1}}) do
    hits
    |> hd()
    |> Map.get("_source")
    |> IO.inspect(label: "object")
    |> JsonHelper.atomize_keys()
  end

  def process_body(body) do
    IO.puts("FULL_BODY")
    IO.inspect(body)
    %{"error" => "unable to process request"}
  end
end
