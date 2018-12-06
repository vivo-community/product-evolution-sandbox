defmodule GraphqlEndpointWeb.Schema.Types do
  use Absinthe.Schema.Notation
  alias GraphqlEndpointWeb.Resolvers

  @desc """
  A person
  """
  object :person do
    field(:uri, :string)
    field(:image, :image)
    field(:name, :name)
    field(:overview_list, list_of(:overview))
    field(:affiliation_list, list_of(:affiliation)) do
      resolve &Resolvers.Affiliations.fetch/3
    end
    field(:education_list, list_of(:education)) do
      resolve &Resolvers.Educations.fetch/3
    end
    field(:grant_list, list_of(:grant)) do
      resolve &Resolvers.Grants.fetch/3
    end
    field(:publication_list, list_of(:publication)) do
      resolve &Resolvers.Publications.fetch/3
    end
  end

  object :image do
    field(:main, :string)
    field(:thumbnail, :string)
  end

  object :name do
    field(:first_name, :string)
    field(:last_name, :string)
    field(:middle_name, :string)
  end

  object :overview do
    field(:overview, :string)
    field(:type, :type)
  end

  object :type do
    field(:code, :string)
    field(:label, :string)
  end

  object :affiliation do
    field(:id, :string)
    field(:label, :string)
    field(:start_date, :date_resolution)
  end

  object :education do
    field(:label, :string)
    field(:org, :organization)
  end

  object :organization do
    field(:id, :string)
    field(:label, :string)
  end

  object :date_resolution do
    field(:date_time, :string)
    field(:resolution, :string)
  end

  #object :funding_role do
    #field(:date_time, :string)
    #field(:label, :string)
  #end

  #object :authorship do
    #field(:date_time, :string)
    #field(:resolution, :string)
  #end

  object :grant do
    field(:id, :string)
    field(:label, :string)
    field(:role_name, :string)
    field(:start_date, :date_resolution)
    field(:end_date, :date_resolution)
  end

  object :venue do
    field(:uri, :string)
    field(:label, :string)
  end

  object :publication do
    field(:id, :string)
    field(:author_list, :string)
    field(:doi, :string)
    field(:label, :string)
    field(:role_name, :string)
    field(:venue, :venue)
  end

end
