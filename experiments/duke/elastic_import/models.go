package widgets_import

import (
	"github.com/jmoiron/sqlx/types"
)

type PersonKeyword struct {
	Uri   string `json:"uri"`
	Label string `json:"label"`
}

type PersonImage struct {
	Main      string `json:"main"`
	Thumbnail string `json:"thumbnail"`
}

type PersonName struct {
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	MiddleName *string `json:"middleName"`
}

type PersonType struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type OverviewType struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type PersonOverview struct {
	Label string       `json:"overview"`
	Type  OverviewType `json:"type"`
}

type Extension struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Person struct {
	Id           string           `json:"id"`
	Uri          string           `json:"uri"`
	SourceId     string           `json:"sourceId"`
	PrimaryTitle string           `json:"primaryTitle"`
	Name         PersonName       `json:"name" elastic:"type:object"`
	Image        PersonImage      `json:"image" elastic:"type:object"`
	Type         PersonType       `json:"type" elastic:"type:object"`
	OverviewList []PersonOverview `json:"overviewList" elastic:"type:nested"`
	KeywordList  []PersonKeyword  `json:"keywordList" elastic:"type:nested"`
	Extensions   []Extension      `json:"extensions" elastic:"type:nested"`
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

type Institution struct {
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
	Id          string      `json:"id"`
	Uri         string      `json:"Uri"`
	Label       string      `json:"label"`
	PersonId    string      `json:"personId"`
	Institution Institution `json:"org" elastic:"type:object"`
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

type Publication struct {
	Id    string `json:"id"`
	Uri   string `json:"uri"`
	Label string `json:"label"`
	// NOTE: this is supposed to be an array
	AuthorList string           `json:"authorList"`
	Doi        string           `json:"doi"`
	Venue      PublicationVenue `json:"venue"`
}

// this is the raw structure in the database
// two json columms:
// * 'data' can be used for change comparison with hash
// * 'data_b' can be used for searches
type Resource struct {
	Uri   string         `db:"uri"`
	Type  string         `db:"type"`
	Hash  string         `db:"hash"`
	Data  types.JSONText `db:"data"`
	DataB types.JSONText `db:"data_b"`
}
