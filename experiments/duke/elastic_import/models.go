package widgets_import

import (
	"github.com/jmoiron/sqlx/types"
	"time"
)

/*
NOTE: model is different per type of keyword

person.keyword
person.keyword.id
person.keyword.uri
person.keyword.label
person.keyword.source
person.keyword.sourceId

publication.keyword
publication.keyword.label
publication.keyword.source
publication.keyword.identifier
*/
type PersonKeyword struct {
	Id       string `json:"id"`
	Uri      string `json:"uri"`
	Label    string `json:"label"`
	Source   string `json:"source"`
	SourceId string `json:"sourceId"`
}

type PublicationKeyword struct {
	Label  string `json:"label"`
	Source string `json:"source"`
	//Identifier string `json:"identifer"`
}

/*
???
type KeywordIdentifier struct {
    Source string `json:"source"`
	SourceId string `json:"sourceId"`
}
*/

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

/*
publication.identifier
publication.identifier.isbn10
publication.identifier.isbn13
publication.identifier.pmid
publication.identifier.doi
publication.identifier.pmcid
*/
type PublicationIdentifier struct {
	Isbn10 string `json:"isbn10"`
	Isbn13 string `json:"isbn13"`
	Pmid   string `json:"pmid"`
	Doi    string `json:"doi"`
	Pmcid  string `json:"pmcid"`
}

// better name for this?
type Type struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

/*
overviewList
overview.id
overview.sourceId
overview.uri
overview.label
overview.type
overview.type.code
overview.type.label
*/
/*
type OverviewType struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}
*/

type PersonOverview struct {
	Label string `json:"label"`
	Type  Type   `json:"type" elastic:"type:object"`
}

/*
serviceRoleList
serviceRole.id
serviceRole.sourceId
serviceRole.uri
serviceRole.label
serviceRole.startDate
serviceRole.startDate.dateTime
serviceRole.startDate.resolution
serviceRole.endDate
serviceRole.endDate.dateTime
serviceRole.endDate.resolution
serviceRole.org
serviceRole.org.id
serviceRole.org.sourceId
serviceRole.org.label
serviceRole.type
serviceRole.type.code
serviceRole.type.label
serviceRole.description
*/
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
	// connector
	PersonId string `json:"personId"`
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

/*
contactList
contact.id
contact.uri
contact.email
contact.email.label
contact.email.type
contact.email.type.code
contact.email.type.label
contact.phone
contact.phone.label
contact.phone.type
contact.phone.type.code
contact.phone.type.label
contact.location
contact.location.label
contact.location.type
contact.location.type.code
contact.location.type.label
contact.website
contact.website.label
contact.website.url
contact.website.type
contact.website.type.code
contact.website.type.label
*/
type Contact struct {
	Id       string   `json:"id"`
	Uri      string   `json:"uri"`
	Email    Email    `json:"email" elastic:"type:object"`
	Phone    Phone    `json:"phone" elastic:"type:object"`
	Location Location `json:"location" elastic:"type:object"`
	Website  Website  `json:"website" elastic:"type:object"`
}

/*
courseTaughtList
courseTaught.id
courseTaught.sourceId
courseTaught.uri
courseTaught.subject
courseTaught.role
courseTaught.courseName
courseTaught.courseNumber
courseTaught.displayDate
courseTaught.startDate
courseTaught.startDate.dateTime
courseTaught.startDate.resolution
courseTaught.endDate
courseTaught.endDate.dateTime
courseTaught.endDate.resolution
courseTaught.org
courseTaught.org.id
courseTaught.org.sourceId
courseTaught.org.uri
courseTaught.org.label
*/
type CourseTaught struct {
	Id           string         `json:"id"`
	SourceId     string         `json:"sourceId"`
	Uri          string         `json:"uri"`
	Subject      string         `json:"subject"`
	Role         string         `json:"role"`
	CourseName   string         `json:"courseName" elastic:"type:object"`
	CourseNumber string         `json":courseNumber" elastic:"type:object"`
	StartDate    DateResolution `json:"startDate" elastic:"type:object"`
	EndDate      DateResolution `json:"endDate" elastic:"type:object"`
	Organization Organization   `json:"organization" elastic:"type:object"`
}

type Extension struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

/*
person
person.id
person.sourceId
person.uri
person.name
person.name.firstName
person.name.middleName
person.name.lastName
person.name.suffix
person.name.prefix
person.identifier
person.identifier.orcid
person.identifier.isni
person.type
person.type.code
person.type.label
person.primaryTitle

person.keyword
person.keyword.id
person.keyword.uri
person.keyword.label
person.keyword.source
person.keyword.sourceId

person.image
person.image.thumbnail
person.image.main
*/
type Person struct {
	Id               string           `json:"id"`
	SourceId         string           `json:"sourceId"`
	Uri              string           `json:"uri"`
	PrimaryTitle     string           `json:"primaryTitle"`
	Name             PersonName       `json:"name" elastic:"type:object"`
	Image            PersonImage      `json:"image" elastic:"type:object"`
	Type             Type             `json:"type" elastic:"type:object"`
	Identifier       PersonIdentifier `json:"identifier" elastic:"type:object"`
	OverviewList     []PersonOverview `json:"overviewList" elastic:"type:nested"`
	KeywordList      []PersonKeyword  `json:"keywordList" elastic:"type:nested"`
	ServiceRoleList  []ServiceRole    `json:"serviceRoleList" elastic:"type:nested"`
	ContactList      []Contact        `json:"contactList" elastic:"type:nested"`
	CourseTaughtList []CourseTaught   `json:"courseTaughtList" elastic:"type:nested"`
	OrganizationList []Organization   `json:"organizationList" elastic:"type:nested"`
	Extensions       []Extension      `json:"extensions" elastic:"type:nested"`
}

type DateResolution struct {
	DateTime   string `json:"dateTime"`
	Resolution string `json:"resolution"`
}

/*
organizationList
organization.id
organization.sourceId
organization.uri
organization.externalUri
organization.label
organization.type
organization.type.label
organization.type.code
organization.childOf
*/
type Organization struct {
	Id          string       `json:"id"`
	SourceId    string       `json:"sourceId"`
	Uri         string       `json:"uri"`
	ExternalUri string       `json:"externalUri"`
	Type        Type         `json:"type" elastic:"type:object"`
	Label       string       `json:"label"`
	// NOTE: can't be recursive, has to be pointer id
	ChildOf     string       `json:"childOf"`
}

/*
type Institution struct {
	Id    string `json:"id"`
	Uri   string `json:"uri"`
	Label string `json:"label"`
}
*/

/*
affiliationRoleList
affiliationRoleList.id
affiliationRoleList.sourceId
affiliationRoleList.uri
affiliationRoleList.label
affiliationRoleList.type
affiliationRoleList.type.code
affiliationRoleList.type.label
affiliationRoleList.startDate
affiliationRoleList.startDate.dateTime
affiliationRoleList.startDate.resolution
affiliationRoleList.endDate
affiliationRoleList.endDate.dateTime
affiliationRoleList.endDate.resolution
affiliationRoleList.organizationID
affiliationRoleList.organizationLabel
*/
type Affiliation struct {
	Id           string         `json:"id"`
	SourceId     string         `json:"sourceId"`
	Uri          string         `json:"uri"`
	Label        string         `json:"label"`
	Type         Type           `json:"type" elastic:"type:object"`
	StartDate    DateResolution `json:"startDate" elastic:"type:object"`
	EndDate      DateResolution `json:"endDate" elastic:"type:object"`
	Organization Organization   `json:"organization" elastic:"type:object"`
	// connector
	PersonId string `json:"personId"`
}

/*
educationList
education.id
education.sourceId
education.uri
education.credential
education.credentialAbbreviation
education.field
education.orgID
education.org
education.org.id
education.org.label
education.dateReceived
education.dateReceived.dateTime
education.dateReceived.resolution
*/
type Education struct {
	Id                     string         `json:"id"`
	SourceId               string         `json:"sourceId"`
	Uri                    string         `json:"uri"`
	Credential             string         `json:"credential"`
	CredentialAbbreviation string         `json:"credentialAbbreviation"`
	Field                  string         `json:"field"`
	Organization           Organization   `json:"org" elastic:"type:object"`
	DateReceived           DateResolution `json:"dateReceived"`
	// connector
	PersonId string `json:"personId"`
}

/*
fundingRoleList
fundingRole.id
fundingRole.sourceId
fundingRole.uri
fundingRole.label
fundingRole.startDate
fundingRole.startDate.dateTime
fundingRole.startDate.resolution
fundingRole.endDate
fundingRole.endDate.dateTime
fundingRole.endDate.resolution
fundingRole.externalId
*/
type FundingRole struct {
	Id  string `json:"id"`
	Uri string `json:"uri"`
	// label for funding role?
	Label string `json:"label"`
	// connectors
	GrantId  string `json:"grantId"`
	PersonId string `json:"personId"`
}

type Grant struct {
	Id         string         `json:"id"`
	SourceId   string         `json:"sourceId"`
	Uri        string         `json:"uri"`
	Label      string         `json:"label"`
	StartDate  DateResolution `json:"startDate" elastic:"type:object"`
	EndDate    DateResolution `json:"endDate" elastic:"type:object"`
	ExternalId string         `json:"externalId"`
}

type Authorship struct {
	Id  string `json:"id"`
	Uri string `json:"uri"`
	Label string `json:"label"`
	// Type Type    `json:"type"`
	// connectors
	PublicationId string `json:"publicationId"`
	PersonId      string `json:"personId"`
}

/*
publication.venue.id
publication.venue.sourceId
publication.venue.uri
publication.venue.label
*/
type PublicationVenue struct {
	Id         string                     `json:"id"`
	SourceId   string                     `json:"sourceId"`
	Uri        string                     `json:"uri"`
	Label      string                     `json:"label"`
	Identifier PublicationVenueIdentifier `json:"identifier" elastic:"type:object"`
}

/*
publication.venue.identifier.eissn
publication.venue.identifier.issn
publication.venue.identifier.isbn10
publication.venue.identifier.isbn13
*/
type PublicationVenueIdentifier struct {
	Isbn10 string `json:"isbn10"`
	Isbn13 string `json:"isbn13"`
	Eissn  string `json:"eissn"`
	Issn   string `json:"issn"`
}

/*
publicationList
publication.id
publication.sourceId
publication.uri
publication.title

publication.identifier
publication.identifier.isbn10
publication.identifier.isbn13
publication.identifier.pmid
publication.identifier.doi
publication.identifier.pmcid

publication.keyword
publication.keyword.id
publication.keyword.sourceId
publication.keyword.uri
publication.keyword.label
publication.keyword.source
publication.keyword.keywordSourceID

publication.venue
publication.venue.id
publication.venue.sourceId
publication.venue.uri
publication.venue.label
publication.venue.identifier.eissn
publication.venue.identifier.issn
publication.venue.identifier.isbn10
publication.venue.identifier.isbn13

publication.dateStandardized
publication.dateStandardized.dateTime
publication.dateStandardized.resolution
publication.dateDisplay
publication.type
publication.type.code
publication.type.label

publication.author
publication.author.blob
publication.author.blob.label
publication.author.blob.personID
publication.author.list
publication.author.list.lastname
publication.author.list.firstname
publication.author.list.rank
publication.author.list.personID

publication.abstract

publication.citations
publication.citations.citationCount
publication.citations.citationDate
publication.citations.citationDate.dateTime
publication.citations.citationDate.resolution
publication.citations.citationSource

publication.pageRange
publication.pageStart
publication.pageEnd
publication.volume
publication.issue

publication.keyword
publication.keyword.label
publication.keyword.source
publication.keyword.identifier

*/
type Publication struct {
	Id         string                `json:"id"`
	SourceId   string                `json:"sourceId"`
	Uri        string                `json:"uri"`
	Title      string                `json:"title"`
	Identifier PublicationIdentifier `json:"identifier" elastic:"type:object"`
	// FIXME: should be actual list?
	AuthorList       string               `json:"authorList"`
	Venue            PublicationVenue     `json:"venue" elastic:"type:object"`
	DateStandardized DateResolution       `json:"dateStandardized" elastic:"type:object"`
	DateDisplay      string               `json:"dateDisplay" elastic:"type:object"`
	Type             Type                 `json:"type" elastic:"type:object"`
	Abstract         string               `json:"abstract"`
	PageRange        string               `json:"pageRange"`
	PageStart        string               `json:"pageStart"`
	PageEnd          string               `json:"pageEnd"`
	Volume           string               `json:"volume"`
	Issue            string               `json:"issue"`
	KeywordList      []PublicationKeyword `json:"keywordList" elastic:"type:nested"`
}

// this is the raw structure in the database
// two json columms:
// * 'data' can be used for change comparison with hash
// * 'data_b' can be used for searches
type Resource struct {
	Uri       string         `db:"uri"`
	Type      string         `db:"type"`
	Hash      string         `db:"hash"`
	Data      types.JSONText `db:"data"`
	DataB     types.JSONText `db:"data_b"`
	CreateAt  time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type StagingResource struct {
	Id       string         `db:"id"`
	Type     string         `db:"type"`
	Data     types.JSONText `db:"data"`
	IsValid  bool           `db:"is_valid"`
	ToDelete bool           `db:"to_delete"`
}
