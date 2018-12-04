package widgets_import

// ********** database json column structs:
// NOTE: this is *not* an independent resource, should it be?
type Keyword struct {
	Uri   string
	Label string
}

// neither is this -in RDF it has to be, but seems like overkill
type DateResolution struct {
	DateTime   string
	Resolution string
}

type ResourceAuthorship struct {
	Uri            string
	PublicationUri string
	PersonUri      string
}

type ResourceFundingRole struct {
	Uri       string
	GrantUri  string
	PersonUri string
	RoleName  string
}

type ResourceGrant struct {
	Uri                      string
	Label                    string
	PrincipalInvestigatorUri string
}

type ResourcePerson struct {
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
	Uri               string
	PersonUri         string
	Label             string
	Start             DateResolution
	OrganizationUri   string
	OrganizationLabel string
}

type ResourceEducation struct {
	Uri       string
	PersonUri string
	Label     string
}

type ResourcePublication struct {
	Uri        string
	Label      string
	AuthorList string
	Doi        string
}

// ********** end database json structs
