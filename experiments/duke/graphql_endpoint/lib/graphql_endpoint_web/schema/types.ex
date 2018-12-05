defmodule GraphqlEndpointWeb.Schema.Types do
  use Absinthe.Schema.Notation

  alias GraphqlEndpointWeb.Resolvers.Common

  @desc """
  A person
  """
  object :person do
    field(:uri, :string)
    field(:image, :image)
    field(:name, :name)
    field(:overview_list, list_of(:overview))
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

  # pe_graphql |   "id" => "per1709582",
  # pe_graphql |   "image" => %{
  # pe_graphql |     "main" => "https://scholars.duke.edu/individual/file_t1709582",
  # pe_graphql |     "thumbnail" => "https://scholars.duke.edu/individual/file_i1709582"
  # pe_graphql |   },
  # pe_graphql |   "keywordList" => nil,
  # pe_graphql |   "name" => %{
  # pe_graphql |     "firstName" => "David",
  # pe_graphql |     "lastName" => "Levi",
  # pe_graphql |     "middleName" => "F."
  # pe_graphql |   },
  # pe_graphql |   "overviewList" => [
  # pe_graphql |     %{
  # pe_graphql |       "overview" => "<p>David F. Levi became the 14th dean of Duke Law School on July 1, 2007. He was named James B. Duke and Benjamin N. Duke Dean of the School of Law in December 2017. Prior to his appointment, he was the Chief United States District Judge for the Eastern District of California with chambers in Sacramento. He was appointed United States Attorney by President Ronald Reagan in 1986 and a United States district judge by President George H. W. Bush in 1990.</p>\r\n<p>A native of Chicago, Dean Levi earned his A.B. in history and literature, magna cum laude, from Harvard College. He entered Harvard's graduate program in history, specializing in English legal history and serving as a teaching fellow in English history and literature. He graduated Order of the Coif in 1980 from Stanford Law School, where he was also president of the Stanford Law Review. Following graduation, he was a law clerk to Judge Ben C. Duniway of the U.S. Court of Appeals for the Ninth Circuit, and then to Justice Lewis F. Powell, Jr., of the U.S. Supreme Court.</p>\r\n<p>Dean Levi has served as chair of two Judicial Conference committees by appointment of the Chief Justice. He was chair of the Civil Rules Advisory Committee (2000-2003) and chair of the Standing Committee on the Rules of Practice and Procedure (2003-2007); he was reappointed to serve as a member of that committee (2009-2015). He was the first president and a founder of the Milton L. Schwartz American Inn of Court, now the Schwartz-Levi American Inn of Court, at the King Hall School of Law, University of California at Davis. He was chair of the Ninth Circuit Task Force on Race, Religious and Ethnic Fairness and was an author of the report of the Task Force. He was president of the Ninth Circuit District Judges Association (2003-2005).</p>\r\n<p>In 2007, Dean Levi was elected a fellow of the American Academy of Arts and Sciences. From 2010 to 2013, he served on the board of directors of Equal Justice Works. In 2014, he was appointed chair of the American Bar Association's Standing Committee on the American Judicial System, and in 2015, he was named co-chair of the North Carolina Commission on the Administration of Law and Justice. He has been elected president of the American Law Institute (ALI), effective May 24, 2017. He is a member of the ALI Council and was an advisor to the ALI's Federal Judicial Code Revision and Aggregate Litigation projects.</p>\r\n<p>Dean Levi is the co-author of Federal Trial Objections (James Publishing 2002). At Duke Law, he has taught courses on judicial behavior, ethics, and legal history.</p>",
  # pe_graphql |       "type" => %{"code" => "overview", "label" => "Overview"}
  # pe_graphql |     }
  # pe_graphql |   ],
  # pe_graphql |   "primaryTitle" => "Professor of Law",
  # pe_graphql |   "sourceId" => "0441800",
  # pe_graphql |   "type" => %{
  # pe_graphql |     "code" => "http://vivoweb.org/ontology/core#FacultyMember",
  # pe_graphql |     "label" => "http://vivoweb.org/ontology/core#FacultyMember"
  # pe_graphql |   },
  # pe_graphql |   "uri" => "https://scholars.duke.edu/individual/per1709582"
end
