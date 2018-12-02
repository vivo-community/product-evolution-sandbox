package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"context"
	// NOTE: empty import is needed cause segfault
	_ "github.com/lib/pq"
	"github.com/olivere/elastic"
	"log"
	"os"
	"text/template"
	"time"
	"flag"
    "github.com/BurntSushi/toml"
)

type Config struct {
  Database database
}

type database struct {
  Server string
  Port int
  Database string
  User string
  Password string
}

type PersonName struct {
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	MiddleName *string `json:"middleName"`
}

type Person struct {
	Uri  string     `json:"uri"`
	Name PersonName `json:"name" elastic:"type:object"`
}

type Mapping struct {
	Definition string
}

type ResourcePerson struct {
	Uri        string
	FirstName  string
	LastName   string
	MiddleName *string
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
		"name":{
			"type":"object",
			"properties": {
				"firstName":  { "type": "text" },
				"lastName":   { "type": "text" },
				"middleName": { "type": "text" }
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

func GetConnection() (*sqlx.DB) {
    return db
}

func listPeople() {
	db = GetConnection()
	resources := []Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1", "Person")
	for _, element := range resources {
		log.Println(element)
		// element is the element from someSlice for where we are
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func makeIndex() {
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

		name := PersonName{resource.FirstName, resource.LastName, resource.MiddleName}
		person := Person{resource.Uri, name}
		put1, err := client.Index().
			Index("people").
			Type("person").
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
	if err != nil{
        log.Println("m=GetPool,msg=connection has failed", err)
    }
	
	makeIndex()
	tryToAdd()
	
	defer db.Close()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
