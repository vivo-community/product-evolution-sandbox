package widgets_import

import (
	"github.com/jmoiron/sqlx/types"
)

// ********** database json column structs:
// NOTE: this is *not* an independent resource, should it be?
type Keyword struct {
	// not sure an 'id' makes sense - they are like #mesh, LOC etc...
	Uri   string
	Label string
}

// neither is this -in RDF it has to be, but seems like overkill
type DateResolution struct {
	DateTime   string
	Resolution string
}

type ResourceFundingRole struct {
	Id       string
	Uri      string
	GrantId  string
	PersonId string
	RoleName string
}

type ResourceGrant struct {
	Id                      string
	Uri                     string
	Label                   string
	PrincipalInvestigatorId string
	Start                   DateResolution
	End                     DateResolution
}

type ResourcePerson struct {
	Id                string
	Uri               string
	AlternateId       string
	FirstName         string
	LastName          string
	MiddleName        *string
	PrimaryTitle      string
	ImageUri          string
	ImageThumbnailUri string
	Type              string
	Overview          string
	Keywords          []Keyword
}

type ResourcePosition struct {
	Id                string
	Uri               string
	PersonId          string
	Label             string
	Start             DateResolution
	OrganizationId    string
	OrganizationLabel string
}

type ResourceInstitution struct {
	Id    string
	Uri   string
	Label string
}

type ResourceEducation struct {
	Id               string
	Uri              string
	PersonId         string
	Label            string
	InsitutionId     string
	InstitutionLabel string
}

type ResourceAuthorship struct {
	Id             string
	Uri            string
	PublicationId  string
	PersonId       string
	AuthorshipType string
}

type ResourcePublication struct {
	Id                  string
	Uri                 string
	Label               string
	AuthorList          string
	Doi                 string
	PublishedIn         string
	PublicationVenueUri string
}

type ResourceOrganization struct {
	Id    string
	Uri   string
	Label string
}

/*** end database json column object maps */

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

// ********** end database json structs
