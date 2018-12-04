package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/OIT-ads-web/widgets_import"
	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"strings"
)

type Resource struct {
	Uri   string         `db:"uri"`
	Type  string         `db:"type"`
	Hash  string         `db:"hash"`
	Data  types.JSONText `db:"data"`
	DataB types.JSONText `db:"data_b"`
}

var client *http.Client

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 50
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
	Attributes struct {
		AuthorList string `json:"authorList"`
		Doi        string `json:"doi"`
	} `json:"attributes"`
}

type Address struct {
	Uri        string `json:"uri"`
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
		PersonUri       string `json:"personUri"`
		DegreeUri       string `json:"degreeUri"`
		Degree          string `json:"degree"`
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
		// could make *string type - or just pass through as ""
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
		FirstName         string  `json:"firstName"`
		LastName          string  `json:"lastName"`
		MiddleName        *string `json:"middleName"`
		PreferredTitle    string  `json:"preferredTitle"`
		PhoneNumber       string  `json:"phoneNumber"`
		PrimaryEmail      string  `json:"primaryEmail"`
		ProfileUrl        string  `json:"profileUrl"`
		ImageUri          string  `json:"imageUri"`
		PrefixName        string  `json:"prefixName"`
		ImageThumbnailUri string  `json:"imageThumbnailUri"`
		AlternateId       string  `json:"alternateId"`
		Overview          string  `json:"overview"`
	} `json:"attributes"`
	Positions     []Position     `json:"positions"`
	Educations    []Education    `json:"educations"`
	Publications  []Publication  `json:"publications"`
	Addresses     []Address      `json:"addresses"`
	ResearchAreas []ResearchArea `json:"researchAreas"`
	Grants        []Grant        `json:"grants"`
}

// ********* end widgets structs

type SolrDoc struct {
	Uri string `json:"URI"`
}

type SolrResults struct {
	Response struct {
		NumFound int       `json:"numFound"`
		Start    int       `json:"start"`
		Docs     []SolrDoc `json:"docs"`
	} `json:"response"`
}

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
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		// TODO: returning 'blank' person, should return nil
		fmt.Println("widgets", err)
		return WidgetsPerson{}
	}

	defer res.Body.Close()

	var person WidgetsPerson
	json.Unmarshal([]byte(body), &person)
	return person
}

//https://stackoverflow.com/questions/2377881/how-to-get-a-md5-hash-from-a-string-in-golang
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
	sqlExists := `SELECT EXISTS (SELECT uri FROM RESOURCES where (uri = $1 AND type =$2))`
	db.Get(&exists, sqlExists, uri, typeName)
	return exists
}

// only add
func addResource(obj interface{}, uri string, typeName string) {
	fmt.Printf("ADD:%v\n", uri)
	db = GetConnection()

	str, err := json.Marshal(obj)
	if err != nil {
		log.Fatalln(err)
	}
	hash := makeHash(string(str))

	res := &Resource{uri, typeName, hash, str, str}

	tx := db.MustBegin()
	sql := `INSERT INTO resources (uri, type, hash, data, data_b) 
	      VALUES (:uri, :type, :hash, :data, :data_b)`
	_, err = tx.NamedExec(sql, res)
	if err != nil {
		log.Fatalln("ERROR(INSERT):%v", err)
	}
	tx.Commit()
}

func saveResource(obj interface{}, uri string, typeName string) {
	db = GetConnection()

	str, err := json.Marshal(obj)
	if err != nil {
		log.Fatalln(err)
	}
	hash := makeHash(string(str))

	found := Resource{}
	res := &Resource{uri, typeName, hash, str, str}

	findSql := `SELECT uri, type, hash, data, data_b  FROM resources 
	  WHERE (uri = $1 AND type = $2)`

	err = db.Get(&found, findSql, uri, typeName)

	tx := db.MustBegin()
	if err != nil {
		fmt.Printf("ADD:%v\n", found.Uri)
		// NOTE: assuming the error means it doesn't exist
		//log.Printf("GET:%v\n", err)
		// must be an add?
		sql := `INSERT INTO resources (uri, type, hash, data, data_b) 
	      VALUES (:uri, :type, :hash, :data, :data_b)`
		_, err := tx.NamedExec(sql, res)
		if err != nil {
			log.Fatalln("ERROR(INSERT):%v", err)
		}
	} else {
		fmt.Printf("UPDATE:%v\n", found.Uri)
		sql := `UPDATE resources 
	    set uri = :uri, 
		type = :type, 
		hash = :hash, 
		data = :data, 
		data_b = :data_b
		WHERE uri = :uri and type = :type`
		_, err := tx.NamedExec(sql, res)

		if err != nil {
			log.Fatalln("ERROR(UPDATE):%v", err)
		}
	}
	tx.Commit()
}

func stashPerson(person WidgetsPerson) {
	fmt.Printf("saving %v\n", person.Uri)
	db = GetConnection()

	// FIXME: if person.Uri is null - should probably exit
	researchAreas := person.ResearchAreas
	var keywords []widgets_import.Keyword
	for _, area := range researchAreas {
		keyword := widgets_import.Keyword{area.Uri, area.Label}
		keywords = append(keywords, keyword)
	}

	obj := widgets_import.ResourcePerson{person.Uri,
		person.Attributes.AlternateId,
		person.Attributes.FirstName,
		person.Attributes.LastName,
		person.Attributes.MiddleName,
		person.Attributes.PreferredTitle,
		person.Attributes.ImageUri,
		person.Attributes.ImageThumbnailUri,
		person.VivoType,
		person.Attributes.Overview,
		keywords}

	saveResource(obj, person.Uri, "Person")
}

func makePositionDate(position Position) widgets_import.DateResolution {
	// NOTE: to make return nil return *DateResolution and ...
	//if position.Attributes.StartYear == nil {
	//	return nil
	//}
	//return &DateResolution{*position.Attributes.StartYear, "year"}
	return widgets_import.DateResolution{position.Attributes.StartYear, "year"}
}

func stashPositions(person WidgetsPerson) {
	fmt.Printf("saving positions:%v\n", person.Uri)
	db = GetConnection()
	positions := person.Positions
	for _, position := range positions {

		start := makePositionDate(position)
		obj := widgets_import.ResourcePosition{position.Uri,
			position.Attributes.PersonUri,
			position.Label,
			start,
			position.Attributes.OrganizationUri,
			position.Attributes.OrganizationLabel}

		saveResource(obj, position.Uri, "Position")
	}
}

func makeIdFromUri(uri string) string {
   return strings.Replace(uri, "https://scholars.duke.edu/individual", "", -1)
}

type Authorship struct {
  PersonId string
  PublicationId string
}

func (auth Authorship) makeUri() string {
   // https://scholars.duke.edu/individual/author1241936-523847
   return fmt.Sprintf("https://scholars.duke.edu/individual/author%s-%s", 
     auth.PublicationId, auth.PersonId)
}

func stashPublications(person WidgetsPerson) {
	fmt.Printf("saving publications:%v\n", person.Uri)
	db = GetConnection()
	publications := person.Publications

	// stash authorships too
	for _, publication := range publications {
		authorship := Authorship{makeIdFromUri(person.Uri), 
		    makeIdFromUri(publication.Uri)} 
		uri := authorship.makeUri()
		fmt.Printf("uri=%v\n", uri)
		//rel := widgets_import.ResourceAuthorship{uri, person.Uri, publication.Uri}
		// TODO: give a new relationship URI
		//saveResource(rel, uri, "Authorship")

		obj := widgets_import.ResourcePublication{publication.Uri,
			publication.Label,
			publication.Attributes.AuthorList,
			publication.Attributes.Doi}
		saveResource(obj, publication.Uri, "Publication")
		//if !resourceExists(publication.Uri, "Publication") {
		//	addResource(obj, publication.Uri, "Publication")
		//}
	}
}

func stashEducations(person WidgetsPerson) {
	fmt.Printf("saving educations:%v\n", person.Uri)
	db = GetConnection()
	educations := person.Educations
	for _, education := range educations {
		obj := widgets_import.ResourceEducation{education.Uri,
			education.Attributes.PersonUri,
			education.Label}
		saveResource(obj, education.Uri, "Education")
	}
}

type FundingRole struct {
	PersonId string
	GrantId string
}

func (role FundingRole) makeUri() string {
  return fmt.Sprintf("http://scholars.duke.edu/individual/investigatorRole%s-%s", 
    role.PersonId, role.GrantId)
}

func stashGrants(person WidgetsPerson) {
	fmt.Printf("saving grants:%v\n", person.Uri)
	db = GetConnection()
	grants := person.Grants

	// stash funding roles too
	for _, grant := range grants {
		// TODO: give a new relationship URI
		fundingRole := FundingRole{makeIdFromUri(person.Uri),
		    makeIdFromUri(grant.Uri)} 
		uri := fundingRole.makeUri()
		fmt.Printf("uri=%v\n", uri)
		//rel := widgets_import.ResourceFundingRole{uri, person.Uri, grant.Uri, grant.Attributes.RoleName}
		// TODO: give a new relationship URI
		//saveResource(rel, uri, "FundingRole")

		obj := widgets_import.ResourceGrant{grant.Uri, grant.Label,
			grant.Attributes.PrincipalInvestigatorUri}
		saveResource(obj, grant.Uri, "Grant")
		//if !resourceExists(grant.Uri, "Grant") {
		//	addResource(obj, grant.Uri, "Grant")
		//}
	}
}

/*** channels ***/
func processUris(cin <-chan string) <-chan WidgetsPerson {
	out := make(chan WidgetsPerson)
	defer wg.Done()
	go func() {
		for line := range cin {
			out <- widgetsParse(line)
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
				case "publications":
					stashPublications(person)
				case "all":
					stashPerson(person)
					stashPositions(person)
					stashEducations(person)
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

func resourceTableExists() bool {
	var exists bool
	db = GetConnection()
	// FIXME: not sure this is right
	sqlExists := `SELECT EXISTS (
        SELECT 1
        FROM   information_schema.tables 
        WHERE  table_catalog = 'vivo_data'
        AND    table_name = 'resources'
    )`
	err := db.QueryRow(sqlExists).Scan(&exists)
	if err != nil {
		log.Fatalln("error checking if row exists %v", err)
	}
	return exists
}

func makeResourceSchema() {
	// NOTE: using data AND data_b columns since binary json
	// does NOT keep ordering, it would mess up
	// any hash based comparison, but it could be still be
	// useful for querying
	sql := `create table resources (
        uri text NOT NULL,
        type text NOT NULL,
        hash text NOT NULL,
        data json NOT NULL,
        data_b jsonb NOT NULL,
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW(),
        PRIMARY KEY(uri, type)
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
	sql := `DELETE from resources`

	switch typeName {
	case "people":
		sql += " WHERE type='Person'"
	case "positions":
		sql += " WHERE type='Position'"
	case "all": // noop
	}
	tx := db.MustBegin()
	tx.MustExec(sql)

	err := tx.Commit()
	if err != nil {
		log.Fatalln("ERROR(DELETE):%v", err)
	}
}

func parseSolr() SolrResults {
	// FIXME: could allow different numbers (for rows) - and/or paging
	// -- 100, 1000 -- NOTE: SolrResults has numFound and start
	//could add-> &sort=timestamp%20asc ??
	url := "https://scholars.duke.edu/vivosolr?q=type:(*FacultyMember)&fl=URI&rows=1000&wt=json"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("widgets", err)
		return SolrResults{}
	}

	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("widgets", err)
		return SolrResults{}
	}

	defer res.Body.Close()

	var results SolrResults
	json.Unmarshal([]byte(body), &results)
	return results

}

func produceUris() <-chan string {
	c := make(chan string)
	defer wg.Done()

	go func() {
		solr := parseSolr()
		for _, doc := range solr.Response.Docs {
			uri := doc.Uri
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
	var configFile string
	flag.StringVar(&configFile, "config", "./config.toml", "a config filename")

	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		fmt.Println("could not find config file, use -c option")
		os.Exit(1)
	}
	fmt.Printf("%#v\n", conf)

	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Database.Server, conf.Database.Port,
		conf.Database.User, conf.Database.Password,
		conf.Database.Database)

	db, err = sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("m=GetPool,msg=connection has failed", err)
	}

	dryRun := flag.Bool("dry-run", false, "just examine widgets parsing")
	typeName := flag.String("type", "people", "type of thing to import")
	remove := flag.Bool("remove", false, "remove existing records")

	flag.Parse()

	if !resourceTableExists() {
		makeResourceSchema()
	}

	// NOTE: either remove OR add?
	if *remove {
		clearResources(*typeName)
	} else {
		wg.Add(3)
		uris := produceUris()
		widgets := processUris(uris)
		persistWidgets(widgets, *dryRun, *typeName)

		wg.Wait()
	}

	defer db.Close()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
