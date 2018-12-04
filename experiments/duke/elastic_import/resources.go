package widgets_import
// ********** database json column structs:
// NOTE: this is *not* an independent resource, should it be?
type Keyword struct {
	Uri   string
	Label string
}

// neither is this -in RDF it has to be, but seems like overkill
type DateResolution struct {
	//Uri        string
	DateTime   string
	Resolution string
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
	Uri string
}

type ResourcePublication struct {
	Uri string
}
// ********** end database json structs


