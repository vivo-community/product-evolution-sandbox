defmodule GraphqlEndpointWeb.Resolvers.People do
  alias GraphqlEndpoint.Search
  alias GraphqlEndpoint.JsonHelper

  def fetch(_parent, %{id: id}, _context) do
    fetch_by_id(id)
  end

  def fetch(_, _, _) do
    {:error, "Unable to retrieve"}
  end

  defp fetch_by_id(id) do
    q = %{query: %{match: %{id: id}}, size: 10}

    case Search.fetch("people", ["person"], q) do
      {:ok, %HTTPoison.Response{status_code: 200, body: body}} ->
        {:ok, process_body(body)}

      {:error, %HTTPoison.Error{reason: reason}} ->
        {:error, reason}
    end
  end

  defp process_body(%{"hits" => _h = %{"hits" => hits}}) do # , "total" => 1
    hits
    |> hd()
    |> Map.get("_source")
    |> IO.inspect(label: "object")
    |> JsonHelper.atomize_understore_keys()
  end

  defp process_body(body) do
    IO.inspect(body, label: "FULL BODY")
    %{"error" => "unable to process request"}
  end
end
