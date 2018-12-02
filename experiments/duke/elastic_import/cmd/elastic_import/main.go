package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	// NOTE: empty import is needed cause segfault
	"flag"
	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq"
	"github.com/olivere/elastic"
	"log"
	"os"
	"text/template"
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

// elastic 'data model'
type PersonKeyword struct {
	Uri   string `json:"uri"`
	Label string `json:"label"`
}

type PersonImage struct {
	Thumbnail string `json:"thumbnail"`
	Main      string `json:"main"`
}

type PersonName struct {
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	MiddleName *string `json:"middleName"`
}

type Person struct {
	Uri          string          `json:"uri"`
	PrimaryTitle string          `json:"primaryTitle"`
	Name         PersonName      `json:"name" elastic:"type:object"`
	Image        PersonImage     `json:"image" elastic:"type:object"`
	KeywordList  []PersonKeyword `json:"keywordList" elastic:"type:nested"`
}
// end elastic data model

// for elastic mapping definitions template
type Mapping struct {
	Definition string
}

// FIXME: centralize - now it's duplicated
// structs to read from resources table
type Keyword struct {
	Uri        string
	Label      string
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

const mappingTemplate = `{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{ 
		{{ .Definition }}
    }
}`

const personMapping = `
"person":{
	"properties":{
		"primaryTitle": { "type": "text" },
		"name":{
			"type":"object",
			"properties": {
				"firstName":  { "type": "text" },
				"lastName":   { "type": "text" },
				"middleName": { "type": "text" }
		    }
		},
		"image": {
			"type": "object",
			"properties": {
				"main":      { "type": "text" },
				"thumbnail": { "type": "text" }
			}
		},
	    "keywordList": {
	      "type": "nested",
	      "properties": {
		      "uri":   { "type": "text" },
		      "label": { "type": "text" }
	      }
	    }
    }
}`

type Resource struct {
	Uri   string         `db:"uri"`
	Type  string         `db:"type"`
	Hash  string         `db:"hash"`
	Data  types.JSONText `db:"data"`
	DataB types.JSONText `db:"data_b"`
}

var psqlInfo string
var db *sqlx.DB

func GetConnection() *sqlx.DB {
	return db
}

func listPeople() {
	db = GetConnection()
	resources := []Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1",
		"Person")
	for _, element := range resources {
		log.Println(element)
		// element is the element from someSlice for where we are
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func listPositions() {
	db = GetConnection()
	resources := []Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1",
		"Position")
	for _, element := range resources {
		log.Println(element)
		// element is the element from someSlice for where we are
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func listPublications() {
	db = GetConnection()
	resources := []Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1",
		"Publication")
	for _, element := range resources {
		log.Println(element)
		// element is the element from someSlice for where we are
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func listEducations() {
	db = GetConnection()
	resources := []Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1",
		"Educations")
	for _, element := range resources {
		log.Println(element)
		// element is the element from someSlice for where we are
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func makePeopleIndex() {
	ctx := context.Background()
	t := template.Must(template.New("index").Parse(mappingTemplate))
	mapping := Mapping{personMapping}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, mapping); err != nil {
		log.Fatalln(err)
	}
	// ??elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	client, err := elastic.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}

	defer client.Stop()

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("people").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("people").BodyString(tpl.String()).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
}

func tryToAdd() {
	ctx := context.Background()

	db = GetConnection()
	resources := []Resource{}

	client, err := elastic.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}

	defer client.Stop()

	err = db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1", "Person")
	for _, element := range resources {
		resource := ResourcePerson{}
		data := element.Data
		json.Unmarshal(data, &resource)

		name := PersonName{resource.FirstName, resource.LastName,
			resource.MiddleName}
		image := PersonImage{resource.ImageUri, resource.ImageThumbnailUri}
		
		var keywordList []PersonKeyword
		for _, keyword := range resource.Keywords {
			pk := PersonKeyword{keyword.Uri, keyword.Label}
		    keywordList = append(keywordList, pk)
		}
		person := Person{resource.Uri, resource.PrimaryTitle, name, image, keywordList}
		put1, err := client.Index().
			Index("people").
			Type("person").
			// TODO: to give ID or not?
			//Id(resource.Uri).
			BodyJson(person).
			Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		fmt.Printf("Indexed person %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		log.Println(element)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

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

	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Database.Server, conf.Database.Port,
		conf.Database.User, conf.Database.Password,
		conf.Database.Database)

	db, err = sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("m=GetPool,msg=connection has failed", err)
	}

	makePeopleIndex()
	tryToAdd()

	defer db.Close()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
