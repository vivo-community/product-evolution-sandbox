package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	// NOTE: empty import is needed
	"context"
	_ "github.com/lib/pq"
	"github.com/olivere/elastic"
	"log"
	"text/template"
	"time"
)

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

const (
	host     = "localhost"
	port     = 5432
	user     = "vivo_data"
	password = "experiment4"
	dbname   = "vivo_data"
)

type Resource struct {
	Uri   string         `db:"uri"`
	Type  string         `db:"type"`
	Hash  string         `db:"hash"`
	Data  types.JSONText `db:"data"`
	DataB types.JSONText `db:"data_b"`
}

var psqlInfo string
var db *sqlx.DB

func init() {
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func listPeople() {
	db, err := sqlx.Connect("postgres", psqlInfo)
	resources := []Resource{}

	err = db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1", "Person")
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
	//err := t.Execute(os.Stdout, r)
	mapping := Mapping{personMapping}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, mapping); err != nil {
		//return err
		log.Fatalln(err)
	}
	//fmt.Println(tpl.String())

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
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

	db, err := sqlx.Connect("postgres", psqlInfo)
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
	defer client.Stop()
}

func main() {
	start := time.Now()
	makeIndex()
	tryToAdd()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
