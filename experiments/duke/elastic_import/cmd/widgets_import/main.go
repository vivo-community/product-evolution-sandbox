package main

import (
	"bufio"
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

type ResourcePerson struct {
	Uri        string
	FirstName  string
	LastName   string
	MiddleName *string
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

func saveAsResource(person WidgetsPerson) {
	fmt.Printf("saving %v\n", person.Uri)
	db = GetConnection()

	//db, err := sqlx.Connect("postgres", psqlInfo)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//defer db.Close()

	// NOTE: if person.Uri is null - should probably exit
	obj := ResourcePerson{person.Uri,
		person.Attributes.FirstName,
		person.Attributes.LastName,
		person.Attributes.MiddleName}

	str, err := json.Marshal(obj)
	if err != nil {
		log.Fatalln(err)
	}
	hash := makeHash(string(str))

	found := Resource{}
	res := &Resource{person.Uri, "Person", hash, str, str}
	fmt.Println(res)

	findSql := `SELECT uri, type, hash, data, data_b  FROM resources 
	  WHERE (uri = $1 AND type = $2)`

	err = db.Get(&found, findSql, person.Uri, "Person")

	tx := db.MustBegin()
	if err != nil {
		// NOTE: assuming the error means it doesn't exist
		log.Println("GET:%v", err)
		// must be an add?
		sql := `INSERT INTO resources (uri, type, hash, data, data_b) 
	      VALUES (:uri, :type, :hash, :data, :data_b)`
		_, err := tx.NamedExec(sql, res)
		if err != nil {
			log.Fatalln("INSERT:%v", err)
		}
	} else {
		fmt.Println("found!!" + found.Uri)
		sql := `UPDATE resources 
	    set uri = :uri, 
		type = :type, 
		hash = :hash, 
		data = :data, 
		data_b = :data_b
		WHERE uri = :uri and type = :type`
		_, err := tx.NamedExec(sql, res)

		if err != nil {
			log.Fatalln("UPDATE:%v", err)
		}
	}
	tx.Commit()
}

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

func persistWidgets(cin <-chan WidgetsPerson, dryRun bool) {
	go func() {
		for person := range cin {
			if dryRun {
				examineParse(person)
			} else {
				saveAsResource(person)
			}
		}
		// 'sink' so need to close waitgroup
		wg.Done()
	}()
}

func produceDuids(filename string) <-chan string {
	c := make(chan string)
	defer wg.Done()

	go func() {
		file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println("could not open file ", filename)
			close(c)
			return
		}
		defer file.Close()

		sc := bufio.NewScanner(file)
		for sc.Scan() {
			c <- sc.Text()
		}
		if err := sc.Err(); err != nil {
			close(c)
			return
		}
		// close
		close(c)
	}()
	return c
}

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

	var filename string

	flag.StringVar(&filename, "f", "", "a filename")
	dryRun := flag.Bool("dry-run", false, "just examine widgets parsing")
	flag.Parse()

	if filename == "" {
		fmt.Println("need -f filename arg")
		os.Exit(1)
	}

	wg.Add(3)
	duids := produceDuids(filename)
	widgets := processDuids(duids)
	persistWidgets(widgets, *dryRun)

	wg.Wait()

	defer db.Close()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
