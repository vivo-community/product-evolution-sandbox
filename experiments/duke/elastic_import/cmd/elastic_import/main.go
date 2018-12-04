package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/OIT-ads-web/widgets_import"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	_ "github.com/lib/pq"
	"github.com/olivere/elastic"
	"log"
	"os"
	"text/template"
	"time"
)

/*
type Config struct {
	Database database
	Elastic  elasticSearch `toml:"elastic"`
}

type elasticSearch struct {
	Url string
}

type database struct {
	Server   string
	Port     int
	Database string
	User     string
	Password string
}
*/

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

type PersonType struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type OverviewType struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type PersonOverview struct {
	// Uri ??
	Label string       `json:"overview"`
	Type  OverviewType `json:"type"`
}

// does overview really need to be a list
type Person struct {
	Uri          string           `json:"uri"`
	SourceId     string           `json:"sourceId"`
	PrimaryTitle string           `json:"primaryTitle"`
	Name         PersonName       `json:"name" elastic:"type:object"`
	Image        PersonImage      `json:"image" elastic:"type:object"`
	Type         PersonType       `json:"type" elastic:"type:object"`
	OverviewList []PersonOverview `json:"overviewList" elastic:"type:nested"`
	KeywordList  []PersonKeyword  `json:"keywordList" elastic:"type:nested"`
}

type Date struct {
	DateTime   string `json:"dateTime"`
	Resolution string `json:"resolution"`
}

type Affiliation struct {
	Uri               string `json:"uri"`
	PersonId          string `json:"personId"`
	Label             string `json:"label"`
	StartDate         Date   `json:"startDate"`
	OrganizationId    string `json:"organizationId"`
	OrganizationLabel string `json:"organizationLabel"`
}

// end elastic data model

// for elastic mapping definitions template
type Mapping struct {
	Definition string
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

// NOTE: dateTime is 'text' because it *can be* nil
const affiliationMapping = `
"affiliation":{
	"properties":{
		"personUri": { "type": "text" },
		"uri":       { "type": "text" },
		"label":     { "type": "text" },
		"startDate": {
			"type": "object",
			"properties": {
				"dateTime":   { "type": "text" },
				"resolution": { "type": "text" }
			}
		},
		"organizationId":   { "type": "text"},
		"organizationLabel": { "type": "text"}
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

func clearIndex(name string) {
	ctx := context.Background()

	client, err := elastic.NewClient(elastic.SetURL(conf.Elastic.Url))
	if err != nil {
		// Handle error
		panic(err)
	}

	defer client.Stop()

	deleteIndex, err := client.DeleteIndex(name).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
	}
}

func clearPeopleIndex() {
	clearIndex("people")
}

func clearAffiliationsIndex() {
	clearIndex("affiliations")
}

func clearResources(typeName string) {
	switch typeName {
	case "people":
		clearPeopleIndex()
	case "affiliations":
		clearAffiliationsIndex()
	case "all":
		clearPeopleIndex()
		clearAffiliationsIndex()
	}
}

// NOTE: 'mappingJson' is just a json string plugged into template
func makeIndex(name string, mappingJson string) {
	ctx := context.Background()
	t := template.Must(template.New("index").Parse(mappingTemplate))
	mapping := Mapping{mappingJson}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, mapping); err != nil {
		log.Fatalln(err)
	}
	client, err := elastic.NewClient(elastic.SetURL(conf.Elastic.Url))
	if err != nil {
		// Handle error
		panic(err)
	}

	defer client.Stop()

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(name).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(name).BodyString(tpl.String()).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
}

func makeAffiliationsIndex() {
	makeIndex("affiliations", affiliationMapping)
}

func makePeopleIndex() {
	makeIndex("people", personMapping)
}

func makeDate(position widgets_import.ResourcePosition) Date {
	// NOTE: to make nullable, return *Date and ...
	//if position.Start == nil {
	//	  return nil
	//}
	//return &Date{position.Start.DateTime, position.Start.Resolution}
	return Date{position.Start.DateTime, position.Start.Resolution}
}

func addAffiliations() {
	ctx := context.Background()

	db = GetConnection()
	resources := []Resource{}

	client, err := elastic.NewClient(elastic.SetURL(conf.Elastic.Url))
	if err != nil {
		// Handle error
		panic(err)
	}

	defer client.Stop()

	err = db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1", "Position")
	for _, element := range resources {
		// NOTE: this is the main difference between types
		resource := widgets_import.ResourcePosition{}
		data := element.Data
		json.Unmarshal(data, &resource)

		// what if blank?
		date := makeDate(resource)

		affiliation := Affiliation{resource.Uri,
			resource.PersonUri,
			resource.Label,
			date,
			resource.OrganizationUri,
			resource.OrganizationLabel}
		put1, err := client.Index().
			Index("affiliations").
			Type("affiliation").
			// TODO: to give ID or not?
			//Id(resource.Uri).
			BodyJson(affiliation).
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

func addPeople() {
	ctx := context.Background()

	db = GetConnection()
	resources := []Resource{}

	client, err := elastic.NewClient(elastic.SetURL(conf.Elastic.Url))
	if err != nil {
		// Handle error
		panic(err)
	}

	defer client.Stop()

	err = db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1", "Person")
	for _, element := range resources {
		// NOTE: this is the main difference between types
		resource := widgets_import.ResourcePerson{}
		data := element.Data
		json.Unmarshal(data, &resource)

		name := PersonName{resource.FirstName, resource.LastName,
			resource.MiddleName}
		image := PersonImage{resource.ImageUri, resource.ImageThumbnailUri}

		personType := PersonType{resource.Type, resource.Type}
		// keywords
		var keywordList []PersonKeyword
		for _, keyword := range resource.Keywords {
			pk := PersonKeyword{keyword.Uri, keyword.Label}
			keywordList = append(keywordList, pk)
		}

		//overviews
		var overviewList []PersonOverview
		overview := PersonOverview{resource.Overview, OverviewType{"overview", "Overview"}}
		overviewList = append(overviewList, overview)
		person := Person{resource.Uri, resource.AlternateId, resource.PrimaryTitle,
			name, image, personType, overviewList, keywordList}
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

func persistResources(dryRun bool, typeName string) {
	if dryRun {
		// what do do here? list?
		//examineParse(person)
		switch typeName {
		case "people":
			listPeople()
		case "affiliations":
			listPositions()
		}
	} else {
		switch typeName {
		case "people":
			makePeopleIndex() /* won't make if already exists */
			addPeople()
		case "affiliations":
			makeAffiliationsIndex()
			addAffiliations()
		}
	}
}

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

	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Database.Server, conf.Database.Port,
		conf.Database.User, conf.Database.Password,
		conf.Database.Database)

	db, err = sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("m=GetPool,msg=connection has failed", err)
	}

	dryRun := flag.Bool("dry-run", false, "just examine resources to be saved")
	remove := flag.Bool("remove", false, "remove existing records")
	typeName := flag.String("type", "people", "type of records to import")

	flag.Parse()

	// NOTE: either remove OR add?
	if *remove {
		clearResources(*typeName)
	} else {
		persistResources(*dryRun, *typeName)
	}

	defer db.Close()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
