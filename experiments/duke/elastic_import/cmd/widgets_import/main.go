package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//"strings"
	"sync"
	"time"
)

type Config struct {
	Database database
}

type database struct {
	Server   string
	Port     int
	Database string
	User     string
	Password string
}

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

type ResearchArea struct {
	Uri        string `json:"uri"`
	Label      string `json:"label"`
	Attributes struct {
		PersonUri string `json:"personUri"`
	} `json:"attributes"`
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
	Label      string `json:"label"`
	Attributes struct {
		PersonUri       string `json:"personUri"`
		DegreeUri       string `json:"degreeUri"`
		Degree          string `json:"degree"`
		OrganizationUri string `json:"organizationUri"`
		Insitution      string `json:"institution"`
	} `json:"attributes"`
}

type Position struct {
	Uri        string `json:"uri"`
	Label      string `json:"label"`
	Attributes struct {
		OrganizationUri   string `json:"organizationUri"`
		OrganizationLabel string `json:"organizationLabel"`
		SchoolUri         string `json:"organizationUri"`
		SchoolLabel       string `json:"organizationLabel"`
		PersonUri         string `json:"personUri"`
	} `json:"attributes"`
}

type WidgetsPerson struct {
	Uri        string `json:"uri"`
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
	} `json:"attributes"`
	Positions     []Position     `json:"positions"`
	Educations    []Education    `json:"educations"`
	Publications  []Publication  `json:"publications"`
	Addresses     []Address      `json:"addresses"`
	ResearchAreas []ResearchArea `json:"researchAreas"`
}

type SolrDoc struct {
	//must remove text vitroIndividual:
	//DocId string `json:"DocId"`
	Uri   string `json:"URI"`
}

type SolrResults struct {
	Response struct {
		NumFound int       `json:"numFound"`
		Start    int       `json:"start"`
		Docs     []SolrDoc `json:"docs"`
	} `json:"response"`
}

func widgetsParse(duid string) WidgetsPerson {
	url := "https://scholars.duke.edu/widgets/api/v0.9/people/complete/all.json?uri=" + duid
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

// this is *not* an independent resource
type Keyword struct {
	Uri   string
	Label string
}

type ResourcePerson struct {
	Uri               string
	FirstName         string
	LastName          string
	MiddleName        *string
	PrimaryTitle      string
	ImageUri          string
	ImageThumbnailUri string
	Keywords          []Keyword
}

type ResourcePosition struct {
	Uri string
}

type ResourceEducation struct {
	Uri string
}

type ResourcePublication struct {
	Uri string
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
	var keywords []Keyword
	for _, area := range researchAreas {
		keyword := Keyword{area.Uri, area.Label}
		keywords = append(keywords, keyword)
	}

	obj := ResourcePerson{person.Uri,
		person.Attributes.FirstName,
		person.Attributes.LastName,
		person.Attributes.MiddleName,
		person.Attributes.PreferredTitle,
		person.Attributes.ImageUri,
		person.Attributes.ImageThumbnailUri,
		keywords}

	saveResource(obj, person.Uri, "Person")
}

func stashPositions(person WidgetsPerson) {
	fmt.Printf("saving positions:%v\n", person.Uri)
	db = GetConnection()
	positions := person.Positions
	for _, position := range positions {
		obj := ResourcePosition{position.Uri}
		saveResource(obj, position.Uri, "Position")
	}
}

func stashPublications(person WidgetsPerson) {
	fmt.Printf("saving publications:%v\n", person.Uri)
	db = GetConnection()
	publications := person.Publications
	for _, publication := range publications {
		obj := ResourcePublication{publication.Uri}
		saveResource(obj, publication.Uri, "Publication")
	}
}

func stashEducations(person WidgetsPerson) {
	fmt.Printf("saving educations:%v\n", person.Uri)
	db = GetConnection()
	educations := person.Educations
	for _, education := range educations {
		obj := ResourceEducation{education.Uri}
		saveResource(obj, education.Uri, "Education")
	}
}

/*** channels ***/
func processDuids(cin <-chan string) <-chan WidgetsPerson {
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

func clearResources() {
	// empty the table first (every time maybe)?
}

func parseSolr() SolrResults {
	// FIXME: could allow different numbers (for rows) - and/or paging
	// -- 100, 1000 ?
	// maybe just faculty?
	//{"numFound":6668,"start":0,
    //could add-> &sort=timestamp%20asc
	// or just URI ?? 
	url := "https://scholars.duke.edu/vivosolr?q=type:(*FacultyMember)&fl=URI&rows=100&wt=json"
	//url := "https://scholars.duke.edu/vivosolr?q=type:(*Person)&fl=DocId&rows=100&wt=json"
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
var conf Config

func main() {
	start := time.Now()
	var err error
	var configFile string
	flag.StringVar(&configFile, "c", "./config.toml", "a config filename")

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
	typeName := flag.String("t", "people", "type of thing to import")

	flag.Parse()

	wg.Add(3)
	uris := produceUris()
	widgets := processDuids(uris)
	persistWidgets(widgets, *dryRun, *typeName)

	wg.Wait()

	defer db.Close()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
