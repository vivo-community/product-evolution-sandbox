defmodule GraphqlEndpointWeb.Resolvers.Grants do
  alias GraphqlEndpoint.Search
  alias GraphqlEndpoint.JsonHelper

  def fetch(_parent, %{person_id: person_id}, _context) do
    fetch_by_person_id(person_id)
  end

  def fetch(parent = %{id: person_id}, _args, _context) do
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

    case Search.fetch("funding-roles", ["funding-role"], q) do
      {:ok, %HTTPoison.Response{status_code: 200, body: body}} ->
        {:ok, process_body(body)}

      {:error, %HTTPoison.Error{reason: reason}} ->
        {:error, reason}
    end
  end

  defp process_body(%{"hits" => _h = %{"hits" => hits}}) do
    Enum.map(hits, &process_funding_role/1)
  end

  defp process_body(body) do
    IO.inspect(body, label: "FULL BODY")
    %{"error" => "unable to process request"}
  end

  defp process_funding_role(es_object) do
    # needs to get underlying grant
    role = Map.get(es_object, "_source")
    grant_id = Map.get(role, "grantId")
    q = %{query: %{match: %{"id" => grant_id}}, size: 100}

    grant =
      case Search.fetch("grants", ["grant"], q) do
        {:ok, %HTTPoison.Response{status_code: 200, body: body}} ->
          {:ok, process_grant_body(body)}

        {:error, %HTTPoison.Error{reason: reason}} ->
          {:error, reason}
      end

    role =
      role
      |> JsonHelper.atomize_understore_keys()

    grant_match = Enum.at(elem(grant, 1), 0)

    compound = %{
      id: grant_match[:id],
      label: grant_match[:label],
      role_name: role[:label],
      start_date: grant_match[:start_date],
      end_date: grant_match[:end_date]
    }

    compound
  end

  defp process_grant_body(%{"hits" => _h = %{"hits" => hits}}) do
    Enum.map(hits, &process_grant/1)
  end

  defp process_grant_body(body) do
    IO.inspect(body, label: "FULL BODY")
    %{"error" => "unable to process request"}
  end

  defp process_grant(es_object) do
    es_object
    |> Map.get("_source")
    |> JsonHelper.atomize_understore_keys()
  end
end
