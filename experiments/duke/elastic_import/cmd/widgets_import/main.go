package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/OIT-ads-web/widgets_import"
	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	"github.com/knakk/rdf"
	_ "github.com/lib/pq"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var client *http.Client

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 100
)

var psqlInfo string

func init() {
	client = createHTTPClient()
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}
	return client
}

// widgets structs
type ResearchArea struct {
	Uri        string `json:"uri"`
	Label      string `json:"label"`
	Attributes struct {
		PersonUri string `json:"personUri"`
	} `json:"attributes"`
}

type Grant struct {
	Uri        string `json:"uri"`
	Label      string `json:"label"`
	VivoType   string `json:"vivoType"`
	Attributes struct {
		PrincipalInvestigatorUri string `json:"piUri"`
		RoleName                 string `json:"roleName"`
		AwardedBy                string `json:"awardedBy"`
		AdministeredBy           string `json:"administeredBy"`
		StartDate                string `json:"startDate"`
		EndDate                  string `json:"endDate"`
	}
}

type Publication struct {
	Uri        string `json:"uri"`
	Label      string `json:"label"`
	VivoType   string `json:"vivoType"`
	Attributes struct {
		// this is kind of like the 'journal'
		PublicationVenue string `json:"publicationVenue"`
		PublishedIn      string `json:"publishedIn"`
		AuthorshipType   string `json:"authorshipType"`
		AuthorList       string `json:"authorList"`
		Volume           string `json:"volume"`
		Issue            string `json:"issue"`
		StartPage        string `json:"startPage"`
		Year             string `json:"year"`
		Source           string `json:"publicationSource"`
		Pmid             string `json:"pmid"`
		Doi              string `json:"doi"`
		Abstract         string `json:"abstract"`
		Datetime         string `json:"datetime"`
	} `json:"attributes"`
}

type Address struct {
	Uri        string `json:"uri"`
	VivoType   string `json:"vivoType"`
	Label      string `json:"label"`
	Attributes struct {
		City       string `json:"city"`
		State      string `json:"state"`
		PostalCode string `json:"postalCode"`
		Address1   string `json:"address1"`
		PersonUri  string `json:"personUri"`
	} `json:"attributes"`
}

type Education struct {
	Uri        string `json:"uri"`
	VivoType   string `json:"vivoType"`
	Label      string `json:"label"`
	Attributes struct {
		PersonUri string `json:"personUri"`
		DegreeUri string `json:"degreeUri"`
		Degree    string `json:"degree"`
		// NOTE: this is not a Duke organization
		OrganizationUri string `json:"organizationUri"`
		Institution     string `json:"institution"`
	} `json:"attributes"`
}

type Position struct {
	Uri        string `json:"uri"`
	Label      string `json:"label"`
	VivoType   string `json:"vivoType"`
	Attributes struct {
		PersonUri         string `json:"personUri"`
		OrganizationUri   string `json:"organizationUri"`
		OrganizationLabel string `json:"organizationLabel"`
		// NOTE: doesn't *always* have school or date
		// so it will send them in as ""
		SchoolUri        string `json:"schoolUri"`
		SchoolLabel      string `json:"schoolLabel"`
		StartDatetimeUri string `json:"startDatetimeUri"`
		StartYear        string `json:"startYear"`
		DateUri          string `json:"dateUri"`
	} `json:"attributes"`
}

type WidgetsPerson struct {
	Uri        string `json:"uri"`
	VivoType   string `json:"vivoType"`
	Attributes struct {
		FirstName              string `json:"firstName"`
		LastName               string `json:"lastName"`
		MiddleName             string `json:"middleName"`
		PreferredTitle         string `json:"preferredTitle"`
		PhoneNumber            string `json:"phoneNumber"`
		PrimaryEmail           string `json:"primaryEmail"`
		ProfileUrl             string `json:"profileUrl"`
		ImageUri               string `json:"imageUri"`
		ImageDownload          string `json:"imageDownload"`
		ImageThumbnailDownload string `json:"imageThumbnailDownload"`
		PrefixName             string `json:"prefixName"`
		ImageThumbnailUri      string `json:"imageThumbnailUri"`
		NetId                  string `json:"netid"`
		AlternateId            string `json:"alternateId"`
		Overview               string `json:"overview"`
	} `json:"attributes"`
	Positions     []Position     `json:"positions"`
	Educations    []Education    `json:"educations"`
	Publications  []Publication  `json:"publications"`
	Addresses     []Address      `json:"addresses"`
	ResearchAreas []ResearchArea `json:"researchAreas"`
	Grants        []Grant        `json:"grants"`
}

// ********* end widgets structs
type WidgetsPersonStub struct {
	Uri string `json:"uri"`
}

type WidgetsOrganization []WidgetsPersonStub

// FIXME should probably return error if fail
func widgetsParse(uri string) WidgetsPerson {
	url := "https://scholars.duke.edu/widgets/api/v0.9/people/complete/all.json?uri=" + uri
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		// TODO: returning a 'blank' person, should return nil
		fmt.Println("widgets", err)
		return WidgetsPerson{}
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("widgets-error", err)
		return WidgetsPerson{}
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		// TODO: returning 'blank' person, should return nil
		fmt.Println("widgets-error", err)
		return WidgetsPerson{}
	}

	defer res.Body.Close()

	var person WidgetsPerson
	json.Unmarshal([]byte(body), &person)
	return person
}

//https://stackoverflow.com/questions/2377881/how-to-get-a-md5-hash-from-a-string-in-golang
func makeHash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

var db *sqlx.DB

func GetConnection() *sqlx.DB {
	return db
}

func examineParse(person WidgetsPerson) {
	fmt.Printf("**********%v\n*************", person.Uri)
	spew.Printf("%+v\n", person)
	fmt.Println("****************")
}

func resourceExists(uri string, typeName string) bool {
	var exists bool
	db = GetConnection()
	sqlExists := `SELECT EXISTS (SELECT id FROM staging where (id = $1 AND type =$2))`
	db.Get(&exists, sqlExists, uri, typeName)
	return exists
}

// only add (presumed existence already checked)
func addResource(obj interface{}, id string, typeName string) {
	fmt.Printf(">ADD:%v\n", id)
	db = GetConnection()

	str, err := json.Marshal(obj)
	if err != nil {
		log.Fatalln(err)
	}

	res := &widgets_import.StagingResource{Id: id, Type: typeName, Data: str}

	tx := db.MustBegin()
	sql := `INSERT INTO STAGING (id, type, data) 
	      VALUES (:id, :type, :data)`
	_, err = tx.NamedExec(sql, res)
	if err != nil {
		log.Fatalln(">ERROR(INSERT):%v", err)
	}
	tx.Commit()
}

func saveResource(obj interface{}, id string, typeName string) {
	db = GetConnection()

	str, err := json.Marshal(obj)
	if err != nil {
		log.Fatalln(err)
	}

	found := widgets_import.StagingResource{}
	res := &widgets_import.StagingResource{Id: id, Type: typeName, Data: str}

	findSql := `SELECT id, type, data FROM staging
	  WHERE (id = $1 AND type = $2)`

	err = db.Get(&found, findSql, id, typeName)

	tx := db.MustBegin()
	if err != nil {
		// NOTE: assuming the error means it doesn't exist
		fmt.Printf(">ADD:%v\n", res.Id)
		sql := `INSERT INTO staging (id, type, data) 
	      VALUES (:id, :type, :data)`
		_, err := tx.NamedExec(sql, res)
		if err != nil {
			log.Fatalln(">ERROR(INSERT):%v", err)
		}
	} else {
		fmt.Printf(">UPDATE:%v\n", found.Id)
		sql := `UPDATE staging
	    set id = id, 
		type = :type, 
		data = :data
		WHERE id = :id and type = :type`
		_, err := tx.NamedExec(sql, res)

		if err != nil {
			log.Fatalln(">ERROR(UPDATE):%v", err)
		}
	}
	tx.Commit()
}

func stashPerson(person WidgetsPerson) {
	fmt.Printf("saving %v\n", person.Uri)
	db = GetConnection()

	// FIXME: if person.Uri is null - should probably exit
	researchAreas := person.ResearchAreas
	var keywords []widgets_import.PersonKeyword
	for _, area := range researchAreas {
		keyword := widgets_import.PersonKeyword{Uri: area.Uri, Label: area.Label}
		keywords = append(keywords, keyword)
	}

	personImage := widgets_import.PersonImage{person.Attributes.ImageDownload,
		person.Attributes.ImageThumbnailDownload}

	// NOTE: this is kind of bogus
	personType := widgets_import.Type{person.VivoType, person.VivoType}
	personName := widgets_import.PersonName{FirstName: person.Attributes.FirstName,
		LastName:   person.Attributes.LastName,
		MiddleName: person.Attributes.MiddleName}

	var overviews []widgets_import.PersonOverview
	overviewType := widgets_import.Type{"overview", "Overview"}
	overview := widgets_import.PersonOverview{person.Attributes.Overview,
		overviewType}
	// NOTE: just an array of one for now
	overviews = append(overviews, overview)

	var extensions []widgets_import.Extension
	extension := widgets_import.Extension{"netid",
		person.Attributes.NetId}
	// NOTE: just an array of one for now
	extensions = append(extensions, extension)

	personId := makeIdFromUri(person.Uri)

	phoneType := widgets_import.Type{Code: "Phone", Label: "Phone"}
	emailType := widgets_import.Type{Code: "Email", Label: "Email"}

	var phones []widgets_import.Phone
	var emails []widgets_import.Email
	var locations []widgets_import.Location
	var websites []widgets_import.Website

	phone := widgets_import.Phone{Label: person.Attributes.PhoneNumber, Type: phoneType}
	email := widgets_import.Email{Label: person.Attributes.PrimaryEmail, Type: emailType}

	phones = append(phones, phone)
	emails = append(emails, email)

	for _, address := range person.Addresses {
		locationType := widgets_import.Type{Code: address.VivoType, Label: address.VivoType}
		location := widgets_import.Location{Label: address.Label, Type: locationType}
		locations = append(locations, location)
	}
	contact := widgets_import.Contact{LocationList: locations, EmailList: emails,
		PhoneList: phones, WebsiteList: websites}

	obj := widgets_import.Person{Id: personId,
		Uri:          person.Uri,
		SourceId:     person.Attributes.AlternateId,
		PrimaryTitle: person.Attributes.PreferredTitle,
		Name:         personName,
		Image:        personImage,
		Contact:      contact,
		Type:         personType,
		OverviewList: overviews,
		KeywordList:  keywords,
		Extensions:   extensions}

	saveResource(obj, personId, "Person")
}

func makePositionDate(position Position) widgets_import.DateResolution {
	return widgets_import.DateResolution{position.Attributes.StartYear, "year"}
}

func stashPositions(person WidgetsPerson) {
	fmt.Printf("saving positions:%v\n", person.Uri)
	db = GetConnection()
	positions := person.Positions
	for _, position := range positions {

		start := makePositionDate(position)
		personId := makeIdFromUri(position.Attributes.PersonUri)
		organizationId := makeIdFromUri(position.Attributes.OrganizationUri)

		org := widgets_import.Organization{Id: organizationId,
			Uri:   position.Attributes.OrganizationUri,
			Label: position.Attributes.OrganizationLabel}

		positionId := makeIdFromUri(position.Uri)
		obj := widgets_import.Affiliation{Id: positionId,
			Uri:          position.Uri,
			PersonId:     personId,
			Label:        position.Label,
			StartDate:    start,
			Organization: org}

		saveResource(obj, positionId, "Position")

		organization := widgets_import.Organization{Id: organizationId,
			Uri:   position.Attributes.OrganizationUri,
			Label: position.Attributes.OrganizationLabel}
		if !resourceExists(organizationId, "Organization") {
			addResource(organization, organizationId, "Organization")
		}
	}
}

func makeIdFromUri(uri string) string {
	return strings.Replace(uri, "https://scholars.duke.edu/individual/", "", -1)
}

func stashEducations(person WidgetsPerson) {
	fmt.Printf("saving educations:%v\n", person.Uri)
	db = GetConnection()
	educations := person.Educations
	for _, education := range educations {
		personId := makeIdFromUri(education.Attributes.PersonUri)

		institutionId := makeIdFromUri(education.Attributes.OrganizationUri)
		institutionUri := education.Attributes.OrganizationUri
		// NOTE: institutions and organizations the same?
		institution := widgets_import.Organization{Id: institutionId,
			Uri:   education.Attributes.OrganizationUri,
			Label: institutionUri}

		educationId := makeIdFromUri(education.Uri)
		obj := widgets_import.Education{Id: educationId,
			Uri: education.Uri,
			// NOTE: these should be different
			Credential:             education.Label,
			CredentialAbbreviation: education.Label,
			PersonId:               personId,
			Organization:           institution}

		saveResource(obj, educationId, "Education")

		if !resourceExists(institutionId, "Institution") {
			addResource(institution, institutionId, "Institution")
		}
	}
}

type FundingRole struct {
	PersonId string
	GrantId  string
}

func (role FundingRole) makeUri() string {
	return fmt.Sprintf("http://scholars.duke.edu/individual/funding-role-%s-%s",
		role.PersonId, role.GrantId)
}

func makeGrantDates(grant Grant) (widgets_import.DateResolution, widgets_import.DateResolution) {
	// NOTE: 'precision' information isn't actually given in widgets data
	start := widgets_import.DateResolution{grant.Attributes.StartDate, "yearMonthDay"}
	end := widgets_import.DateResolution{grant.Attributes.EndDate, "yearMonthDay"}
	return start, end
}

func stashGrants(person WidgetsPerson) {
	fmt.Printf("saving grants:%v\n", person.Uri)
	db = GetConnection()
	grants := person.Grants

	// NOTE: stashes funding roles AND grants
	for _, grant := range grants {
		personId := makeIdFromUri(person.Uri)
		grantId := makeIdFromUri(grant.Uri)
		fundingRoleId := fmt.Sprintf("%s-%s", grantId, personId)
		fundingRole := FundingRole{personId, grantId}

		// NOTE: this is an approximation of real function, uri is fake
		uri := fundingRole.makeUri()
		rel := widgets_import.FundingRole{Id: fundingRoleId,
			Uri:      uri,
			GrantId:  grantId,
			PersonId: personId,
			Label:    grant.Attributes.RoleName}
		saveResource(rel, fundingRoleId, "FundingRole")

		//pi := makeIdFromUri(grant.Attributes.PrincipalInvestigatorUri)
		start, end := makeGrantDates(grant)
		obj := widgets_import.Grant{Id: grantId,
			Uri:       grant.Uri,
			Label:     grant.Label,
			StartDate: start,
			EndDate:   end}
		if !resourceExists(grantId, "Grant") {
			addResource(obj, grantId, "Grant")
		}
	}
}

// NOTE: just an intermediary object
type Authorship struct {
	PublicationId string
	PersonId      string
}

func (auth Authorship) makeUri() string {
	return fmt.Sprintf("https://scholars.duke.edu/individual/authorship-%s-%s",
		auth.PublicationId, auth.PersonId)
}

func stashPublications(person WidgetsPerson) {
	fmt.Printf("saving publications:%v\n", person.Uri)
	db = GetConnection()
	publications := person.Publications

	for _, publication := range publications {
		personId := makeIdFromUri(person.Uri)
		publicationId := makeIdFromUri(publication.Uri)
		authorshipId := fmt.Sprintf("%s-%s", publicationId, personId)
		authorship := Authorship{publicationId, personId}

		uri := authorship.makeUri()
		rel := widgets_import.Authorship{Id: authorshipId,
			Uri:           uri,
			PublicationId: publicationId,
			PersonId:      personId,
			// FIXME: should maybe be 'Type'
			Label: publication.Attributes.AuthorshipType}
		saveResource(rel, authorshipId, "Authorship")

		venue := widgets_import.PublicationVenue{
			Uri:   publication.Attributes.PublicationVenue,
			Label: publication.Attributes.PublishedIn}

		identifier := widgets_import.PublicationIdentifier{Doi: publication.Attributes.Doi}
		obj := widgets_import.Publication{Id: publicationId,
			Uri:         publication.Uri,
			Title:       publication.Label,
			Abstract:    publication.Attributes.Abstract,
			DateDisplay: publication.Attributes.Datetime,
			Issue:       publication.Attributes.Issue,
			PageStart:   publication.Attributes.StartPage,
			AuthorList:  publication.Attributes.AuthorList,
			Identifier:  identifier,
			Venue:       venue}
		if !resourceExists(publicationId, "Publication") {
			addResource(obj, publicationId, "Publication")
		}
	}
}

func processUri(uri string) WidgetsPerson {
	person := widgetsParse(uri)
	return person
}

func persistPerson(person WidgetsPerson, dryRun bool, typeName string) {
	if dryRun {
		examineParse(person)
	} else {
		switch typeName {
		case "people":
			stashPerson(person)
		case "positions":
			stashPositions(person)
		case "educations":
			stashEducations(person)
		case "grants":
			stashGrants(person)
		case "publications":
			stashPublications(person)
		case "all":
			stashPerson(person)
			stashPositions(person)
			stashEducations(person)
			stashGrants(person)
			stashPublications(person)
		default:
			stashPerson(person)
		}
	}
}

/*** channels ***/
func processUris(cin <-chan string) <-chan WidgetsPerson {
	out := make(chan WidgetsPerson)
	defer wg.Done()
	go func() {
		for line := range cin {
			person := widgetsParse(line)
			if person.Uri != "" {
				//out <- widgetsParse(line)
				out <- person
			}
		}
		close(out)
	}()
	return out
}

func persistWidgets(cin <-chan WidgetsPerson, dryRun bool, typeName string) {
	go func() {
		for person := range cin {
			if dryRun {
				examineParse(person)
			} else {
				switch typeName {
				case "people":
					stashPerson(person)
				case "positions":
					stashPositions(person)
				case "educations":
					stashEducations(person)
				case "grants":
					stashGrants(person)
				case "publications":
					stashPublications(person)
				case "all":
					stashPerson(person)
					stashPositions(person)
					stashEducations(person)
					stashGrants(person)
					stashPublications(person)
				default:
					stashPerson(person)
				}
			}
		}
		// 'sink' so need to close waitgroup
		wg.Done()
	}()
}

func stagingTableExists() bool {
	var exists bool
	db = GetConnection()
	// FIXME: not sure this is right
	sqlExists := `SELECT EXISTS (
        SELECT 1
        FROM   information_schema.tables 
        WHERE  table_catalog = 'vivo_data'
        AND    table_name = 'staging'
    )`
	err := db.QueryRow(sqlExists).Scan(&exists)
	if err != nil {
		log.Fatalln("error checking if row exists %v", err)
	}
	return exists
}

// 'type' should match up to a schema
func makeStagingSchema() {
	sql := `create table staging (
        id text NOT NULL,
        type text NOT NULL,
        data json NOT NULL,
		is_valid boolean DEFAULT FALSE,
		to_delete boolean DEFAULT FALSE,
        PRIMARY KEY(id, type)
    )`

	db = GetConnection()
	tx := db.MustBegin()
	tx.MustExec(sql)

	err := tx.Commit()
	if err != nil {
		log.Fatalln("ERROR(CREATE):%v", err)
	}
}

func clearResources(typeName string) {
	db = GetConnection()
	sql := `DELETE from staging`

	switch typeName {
	case "people":
		sql += " WHERE type='Person'"
	case "positions":
		// NOTE: organization only come from Positions (now)
		sql += " WHERE type='Position' or type ='Organization'"
	case "grants":
		sql += " WHERE type='Grant' or type='FundingRole'"
	case "publications":
		sql += " WHERE type='Publication' or type='Authorship'"
	case "educations":
		// NOTE: institutions only come from Educations (now)
		sql += " WHERE type='Education' or type='Institution'"
	case "all": // noop
	}
	tx := db.MustBegin()
	tx.MustExec(sql)

	log.Println(sql)
	err := tx.Commit()
	if err != nil {
		log.Fatalln(">ERROR(DELETE):%v", err)
	}
}

// examples:
// * computer science=https://scholars.duke.edu/individual/org50000500
// * trinity=https://scholars.duke.edu/individual/org50000491
func parseOrganizationPage(orgUri string) WidgetsOrganization {
	url := "https://scholars.duke.edu/widgets/api/v0.9/organizations/people/5.json?uri=" + orgUri
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("widgets", err)
		return WidgetsOrganization{}
	}

	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("widgets", err)
		return WidgetsOrganization{}
	}

	defer res.Body.Close()

	var results WidgetsOrganization
	json.Unmarshal([]byte(body), &results)
	return results
}

// from Org
func gatherUrisFromWidgetsOrg(org *string) []string {
	var uris []string
	orgs := parseOrganizationPage("https://scholars.duke.edu/individual/" + *org)
	for _, doc := range orgs {
		uri := doc.Uri
		fmt.Println(uri)
		uris = append(uris, uri)
	}
	return uris
}

// from Rdf file
func gatherUrisFromRdfFile(fileName string) []string {
	var uris []string
	f, err := os.Open(fileName)
	if err != nil {
		// handle error
	}
	dec := rdf.NewTripleDecoder(f, rdf.RDFXML)
	for triple, err := dec.Decode(); err != io.EOF; triple, err = dec.Decode() {
		// do something with triple ..
		fmt.Println(triple.Subj)
		uris = append(uris, triple.Subj.String())
	}
	return uris
}

// 3 hrs for medicine
func produceUrisFromWidgetsOrg(org *string) <-chan string {
	c := make(chan string)
	defer wg.Done()

	go func() {
		org := parseOrganizationPage("https://scholars.duke.edu/individual/" + *org)
		for _, doc := range org {
			uri := doc.Uri
			c <- uri
		}
		close(c)
	}()
	return c
}

func produceUrisFromRdfFile(fileName string) <-chan string {
	c := make(chan string)
	defer wg.Done()
	f, err := os.Open(fileName)
	if err != nil {
		// handle error
	}
	dec := rdf.NewTripleDecoder(f, rdf.RDFXML)
	go func() {
		for triple, err := dec.Decode(); err != io.EOF; triple, err = dec.Decode() {
			uri := triple.Subj.String()
			c <- uri
		}
		close(c)
	}()
	return c
}

func readPeopleList() (io.Reader, error) {
	//https://scholars.duke.edu/listrdf?vclass=http%3A%2F%2Fxmlns.com%2Ffoaf%2F0.1%2FPerson
	url := "https://scholars.duke.edu/listrdf?vclass=http://xmlns.com/foaf/0.1/Person"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/rdf+xml")

	if err != nil {
		fmt.Println("ERROR:", err)
	}

	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("ERROR:", err)
	}

	defer res.Body.Close()

	reader := strings.NewReader(string(body))
	return reader, err
	//var results WidgetsOrganization
	//json.Unmarshal([]byte(body), &results)
	//return results
}

func produceUrisFromVivo() <-chan string {
	c := make(chan string)
	defer wg.Done()
	f, err := readPeopleList()

	if err != nil {
		// handle error
	}
	dec := rdf.NewTripleDecoder(f, rdf.RDFXML)
	go func() {
		for triple, err := dec.Decode(); err != io.EOF; triple, err = dec.Decode() {
			uri := triple.Subj.String()
			c <- uri
		}
		close(c)
	}()
	return c
}

/**** end channels ****/
var wg sync.WaitGroup
var conf widgets_import.Config

func main() {
	start := time.Now()
	var err error
	var rdfFile string

	flag.StringVar(&rdfFile, "rdf", "", "an rdf file (of person uris)")

	if os.Getenv("ENVIRONMENT") == "development" {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")

		value, exists := os.LookupEnv("CONFIG_PATH")
		if exists {
			viper.AddConfigPath(value)
		}

		viper.ReadInConfig()
	} else {
		replacer := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(replacer)
		viper.BindEnv("database.server")
		viper.BindEnv("database.port")
		viper.BindEnv("database.database")
		viper.BindEnv("database.user")
		viper.BindEnv("database.password")
		viper.BindEnv("elastic.url")
	}

	dryRun := flag.Bool("dry-run", false, "just examine widgets parsing")
	typeName := flag.String("type", "people", "type of thing to import")
	source := flag.String("source", "org", "source of data")
	remove := flag.Bool("remove", false, "remove existing records")
	org := flag.String("org", "org50000500", "which org id to import (defaults to CS)")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if err := viper.Unmarshal(&conf); err != nil {
		fmt.Printf("could not establish read into conf structure %s\n", err)
		os.Exit(1)
	}

	//viper.SetDefault("elastic.url", "http://localhost:9200")

	fmt.Printf("trying to connect to %s\n", conf.Database.Server)
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Database.Server, conf.Database.Port,
		conf.Database.User, conf.Database.Password,
		conf.Database.Database)

	db, err = sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("m=GetPool,msg=connection has failed", err)
	}

	if !stagingTableExists() {
		makeStagingSchema()
	}

	// NOTE: either remove OR add?
	if *remove {
		clearResources(*typeName)
	} else {
		// 1. this way
		wg.Add(3)
		var uris <-chan string
		if len(rdfFile) > 0 {
			uris = produceUrisFromRdfFile(rdfFile)
		} else {
			switch *source {
			// default
			case "org":
				uris = produceUrisFromWidgetsOrg(org)
			case "vivo":
				uris = produceUrisFromVivo()
			}
		}

		widgets := processUris(uris)
		persistWidgets(widgets, *dryRun, *typeName)
		wg.Wait()
	}

	defer db.Close()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
