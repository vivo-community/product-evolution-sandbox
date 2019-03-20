package graphql_endpoint

// these are elastic json models
type Type struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type PersonKeyword struct {
	Uri   string `json:"uri"`
	Label string `json:"label"`
}

type PersonImage struct {
	Main      string `json:"main"`
	Thumbnail string `json:"thumbnail"`
}

type PersonName struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	MiddleName string `json:"middleName"`
	Suffix     string `json:"suffix"`
	Prefix     string `json:"prefix"`
}

type PersonIdentifier struct {
	Orcid string `json:"orchid"`
	Isni  string `json:"isni"`
}

type OverviewType struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type PersonOverview struct {
	Label string       `json:"overview"`
	Type  OverviewType `json:"type"`
}

type ServiceRole struct {
	Id           string         `json:"id"`
	SourceId     string         `json:"sourceId"`
	Uri          string         `json:"uri"`
	Label        string         `json:"label"`
	Description  string         `json:"description"`
	StartDate    DateResolution `json:"startDate" elastic:"type:object"`
	EndDate      DateResolution `json:"endDate" elastic:"type:object"`
	Organization Organization   `json:"organization" elastic:"type:object"`
	Type         Type           `json:"type" elastic:"type:object"`
	PersonId     string         `json:"personId"`
}

type Email struct {
	Label string `json:"label"`
	Type  Type   `json:"type" elastic:"type:object"`
}

type Phone struct {
	Label string `json:"label"`
	Type  Type   `json:"type" elastic:"type:object"`
}

type Location struct {
	Label string `json:"label"`
	Type  Type   `json:"type" elastic:"type:object"`
}

type Website struct {
	Label string `json:"label"`
	Url   string `json:"url"`
	Type  Type   `json:"type" elastic:"type:object"`
}

type Contact struct {
	EmailList    []Email    `json:"emailList" elastic:"type:nested"`
	PhoneList    []Phone    `json:"phoneList" elastic:"type:nested"`
	LocationList []Location `json:"locationList" elastic:"type:nested"`
	WebsiteList  []Website  `json:"websiteList" elastic:"type:nested"`
}

type CourseTaught struct {
	Id           string         `json:"id"`
	SourceId     string         `json:"sourceId"`
	Uri          string         `json:"uri"`
	Subject      string         `json:"subject"`
	Role         string         `json:"role"`
	CourseName   string         `json:"courseName" elastic:"type:object"`
	CourseNumber string         `json:"courseNumber" elastic:"type:object"`
	StartDate    DateResolution `json:"startDate" elastic:"type:object"`
	EndDate      DateResolution `json:"endDate" elastic:"type:object"`
	Organization Organization   `json:"organization" elastic:"type:object"`
}

type Extension struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Person struct {
	Id               string           `json:"id"`
	Uri              string           `json:"uri"`
	SourceId         string           `json:"sourceId"`
	PrimaryTitle     string           `json:"primaryTitle"`
	Name             PersonName       `json:"name" elastic:"type:object"`
	Image            PersonImage      `json:"image" elastic:"type:object"`
	Type             Type             `json:"type" elastic:"type:object"`
	Identifier       PersonIdentifier `json:"identifier" elastic:"type:object"`
	Contact          Contact          `json:"contact" elastic:"type:object"`
	OverviewList     []PersonOverview `json:"overviewList" elastic:"type:nested"`
	KeywordList      []PersonKeyword  `json:"keywordList" elastic:"type:nested"`
	ServiceRoleList  []ServiceRole    `json:"serviceRoleList" elastic:"type:nested"`
	CourseTaughtList []CourseTaught   `json:"courseTaughtList" elastic:"type:nested"`
	Extensions       []Extension      `json:"extensions" elastic:"type:nested"`
	EducationList    []Education      `json:"educationList" elastic:"type:nested"`
	AffliationList   []Affiliation    `json:"affiliationList" elastic:"type:nested"`
}

type DateResolution struct {
	DateTime   string `json:"dateTime"`
	Resolution string `json:"resolution"`
}

type Organization struct {
	Id    string `json:"id"`
	Uri   string `json:"uri"`
	Label string `json:"label"`
}

type Affiliation struct {
	Id           string         `json:"id"`
	Uri          string         `json:"uri"`
	PersonId     string         `json:"personId"`
	Label        string         `json:"label"`
	StartDate    DateResolution `json:"startDate"`
	Organization Organization   `json:"organization"`
}

type Education struct {
	Id                     string         `json:"id"`
	Uri                    string         `json:"uri"`
	Credential             string         `json:"credential"`
	CredentialAbbreviation string         `json:"credentialAbbreviation"`
	DateRecieved           DateResolution `json:"dateReceived"`
	PersonId               string         `json:"personId"`
	Organization           Organization   `json:"organization" elastic:"type:object"`
}

type FundingRole struct {
	Id       string `json:"id"`
	Uri      string `json:"uri"`
	GrantId  string `json:"grantId"`
	PersonId string `json:"personId"`
	Label    string `json:"label"`
}

type Grant struct {
	Id        string         `json:"id"`
	Uri       string         `json:"uri"`
	Label     string         `json:"label"`
	StartDate DateResolution `json:"startDate"`
	EndDate   DateResolution `json:"endDate"`
}

type Authorship struct {
	Id            string `json:"id"`
	Uri           string `json:"uri"`
	PublicationId string `json:"publicationId"`
	PersonId      string `json:"personId"`
	Label         string `json:"label"`
}

type PublicationVenue struct {
	Uri   string `json:"uri"`
	Label string `json:"label"`
}

type PublicationIdentifier struct {
	Isbn10 string `json:"isbn10"`
	Isbn13 string `json:"isbn13"`
	Pmid   string `json:"pmid"`
	Doi    string `json:"doi"`
	Pmcid  string `json:"pmcid"`
}

type PublicationKeyword struct {
	Label  string `json:"label"`
	Source string `json:"source"`
}

type Publication struct {
	Id       string `json:"id"`
	SourceId string `json:"sourceId"`
	Uri      string `json:"uri"`
	Title    string `json:"title"`
	// NOTE: this is supposed to be an array
	AuthorList       string                `json:"authorList"`
	Venue            PublicationVenue      `json:"venue"`
	Identifier       PublicationIdentifier `json:"identifier"`
	DateStandardized DateResolution        `json:"dateStandardized"`
	DateDisplay      string                `json:"dateDisplay"`
	Type             Type                  `json:"type"`
	Abstract         string                `json:"abstract"`
	PageRange        string                `json:"pageRange"`
	PageStart        string                `json:"pageStart"`
	PageEnd          string                `json:"pageEnd"`
	Volume           string                `json:"volume"`
	Issue            string                `json:"issue"`
	KeywordList      []PublicationKeyword  `json:"keywordList" elastic:"type:nested"`
}

// these are graphql only json models
type PageInfo struct {
	PerPage     int `json:"perPage"`
	CurrentPage int `json:"page"`
	TotalPages  int `json:"totalPages"`
	Count       int `json:"count"`
}

type Facet struct {
	Label string `json:"label"`
	Count int64  `json:"count"`
}

type PeopleFacets struct {
	Keywords    []Facet `json:"keywords"`
	Types       []Facet `json:"types"`
	Departments []Facet `json:"departments"`
}

type PersonList struct {
	Results  []Person      `json:"data"`
	PageInfo PageInfo      `json:"pageInfo"`
	Facets   *PeopleFacets `json:"facets"`
}

type PublicationList struct {
	Results  []Publication `json:"data"`
	PageInfo PageInfo      `json:"pageInfo"`
}

type GrantList struct {
	Results  []Grant  `json:"data"`
	PageInfo PageInfo `json:"pageInfo"`
}
