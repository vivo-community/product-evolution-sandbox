package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/OIT-ads-web/widgets_import"
	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/olivere/elastic"
)

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
		"id":           { "type": "text" },
		"uri":          { "type": "text" },
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
		},
		"extensions": {
			"type": "nested",
			"properties": {
				"key":   { "type": "text" },
				"value": { "type": "text" }
			}			
		}
    }
}`

// NOTE: dateTime is 'text' because it *can be* nil
const affiliationMapping = `
"affiliation":{
	"properties":{
		"id":        { "type": "text" },
		"uri":       { "type": "text" },
		"personId":  { "type": "text" },
		"label":     { "type": "text" },
		"startDate": {
			"type": "object",
			"properties": {
				"dateTime":   { "type": "text" },
				"resolution": { "type": "text" }
			}
		},
		"organizationId":    { "type": "text" },
		"organizationLabel": { "type": "text" } 
    }
}`

const educationMapping = `
"education":{
	"properties":{
		"id":        { "type": "text" },
		"uri":       { "type": "text" },
		"label":     { "type": "text" },
		"personId":  { "type": "text" },
		"org":     { 
			"type": "object",
			"properties": {
				"id": { "type": "text" },
				"label": { "type": "text" }
			}
		}
	}
}`

const grantMapping = `
"grant":{
	"properties":{
		"id":        { "type": "text" },
		"uri":       { "type": "text" },
		"label":     { "type": "text" },
		"startDate": {
			"type": "object",
			"properties": {
				"dateTime":   { "type": "text" },
				"resolution": { "type": "text" }
			}
		},
		"endDate": {
			"type": "object",
			"properties": {
				"dateTime":   { "type": "text" },
				"resolution": { "type": "text" }
			}
		}
	}
}`

const fundingRoleMapping = `
"funding-role":{
	"properties":{
		"id":        { "type": "text" },
		"uri":       { "type": "text" },
		"grantId":   { "type": "text" },
		"personId":  { "type": "text" },
		"label":     { "type": "text" }
	}
}`

const publicationMapping = `
"publication":{
	"properties":{
		"id":         { "type": "text" },
		"uri":        { "type": "text" },
		"label":      { "type": "text" },
		"authorList": { "type": "text" },
		"doi":        { "type": "text" },
        "venue":      { 
			"type": "object",
			"properties": {
				"uri":   { "type": "text" },
				"label": { "type": "text" }
			}
		}
	}
}`

const authorshipMapping = `
"authorship":{
	"properties":{
		"id":             { "type": "text" },
		"uri":            { "type": "text" },
		"publicationId":  { "type": "text" },
		"personId":       { "type": "text" },
		"label":          { "type": "text" }
	}
}`

var psqlInfo string
var db *sqlx.DB
var client *elastic.Client

func GetConnection() *sqlx.DB {
	return db
}

func GetClient() *elastic.Client {
	return client
}

func retrieveType(typeName string) []widgets_import.Resource {
	db = GetConnection()
	resources := []widgets_import.Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data FROM resources WHERE type =  $1", typeName)
	if err != nil {
		log.Fatalln(err)
	}
	return resources
}

func listType(typeName string) {
	db = GetConnection()
	resources := []widgets_import.Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data FROM resources WHERE type =  $1",
		typeName)
	for _, element := range resources {
		log.Println(element)
		// element is the element from someSlice for where we are
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func listPeople() {
	listType("Person")
}

func listPositions() {
	listType("Position")
}

func listEducations() {
	listType("Education")
}

func listGrants() {
	listType("Grant")
}

func listFundingRoles() {
	listType("FundingRole")
}

func listPublications() {
	listType("Publication")
}

func clearIndex(name string) {
	ctx := context.Background()

	client = GetClient()

	deleteIndex, err := client.DeleteIndex(name).Do(ctx)
	if err != nil {
		log.Printf("ERROR:%v\n", err)
		return
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
		log.Println("Not acknowledged")
	} else {
		log.Println("Acknowledged!")
	}
}

func clearPeopleIndex() {
	clearIndex("people")
}

func clearAffiliationsIndex() {
	clearIndex("affiliations")
}

func clearEducationsIndex() {
	clearIndex("educations")
}

func clearGrantsIndex() {
	clearIndex("grants")
}

func clearFundingRolesIndex() {
	clearIndex("funding-roles")
}

func clearPublicationsIndex() {
	clearIndex("publications")
}

func clearAuthorshipsIndex() {
	clearIndex("authorships")
}

func clearResources(typeName string) {
	switch typeName {
	case "people":
		clearPeopleIndex()
	case "affiliations":
		clearAffiliationsIndex()
	case "educations":
		clearEducationsIndex()
	case "grants":
		clearGrantsIndex()
		clearFundingRolesIndex()
	case "publications":
		clearPublicationsIndex()
		clearAuthorshipsIndex()
	case "all":
		clearPeopleIndex()
		clearAffiliationsIndex()
		clearEducationsIndex()
		clearGrantsIndex()
		clearFundingRolesIndex()
		clearPublicationsIndex()
		clearAuthorshipsIndex()
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
	client = GetClient()

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

func makePeopleIndex() {
	makeIndex("people", personMapping)
}

func makeAffiliationsIndex() {
	makeIndex("affiliations", affiliationMapping)
}

func makeEducationsIndex() {
	makeIndex("educations", educationMapping)
}

func makeGrantsIndex() {
	makeIndex("grants", grantMapping)
}

func makeFundingRolesIndex() {
	makeIndex("funding-roles", fundingRoleMapping)
}

func makePublicationsIndex() {
	makeIndex("publications", publicationMapping)
}

func makeAuthorshipsIndex() {
	makeIndex("authorships", authorshipMapping)
}

// TODO: sketch of way to make slightly more generic
// every type will have different db -> elastic mapping though
func addToIndex(index string, typeName string, obj interface{}) {
	ctx := context.Background()
	client = GetClient()

	put1, err := client.Index().
		Index(index).
		Type(typeName).
		BodyJson(obj).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	spew.Println(obj)
}

func addPeople() {
	people := retrieveType("Person")
	for _, element := range people {
		resource := widgets_import.Person{}
		data := element.Data
		json.Unmarshal(data, &resource)

		addToIndex("people", "person", resource)
	}
}

func addAffiliations() {
	positions := retrieveType("Position")
	for _, element := range positions {
		resource := widgets_import.Affiliation{}
		data := element.Data
		json.Unmarshal(data, &resource)

		addToIndex("affiliations", "affiliation", resource)
	}
}

func addEducations() {
	educations := retrieveType("Education")
	for _, element := range educations {
		resource := widgets_import.Education{}
		data := element.Data
		json.Unmarshal(data, &resource)

		addToIndex("educations", "education", resource)
	}
}

func addGrants() {
	grants := retrieveType("Grant")
	for _, element := range grants {
		resource := widgets_import.Grant{}
		data := element.Data
		json.Unmarshal(data, &resource)

		addToIndex("grants", "grant", resource)
	}
}

func addFundingRoles() {
	fundingRoles := retrieveType("Funding-Role")
	for _, element := range fundingRoles {
		resource := widgets_import.FundingRole{}
		data := element.Data
		json.Unmarshal(data, &resource)

		addToIndex("funding-roles", "funding-role", resource)
	}
}

func addPublications() {
	publications := retrieveType("Publication")
	for _, element := range publications {
		resource := widgets_import.Publication{}
		data := element.Data
		json.Unmarshal(data, &resource)

		addToIndex("publications", "publication", resource)
	}
}

func addAuthorships() {
	authorships := retrieveType("Authorship")
	for _, element := range authorships {
		resource := widgets_import.Authorship{}
		data := element.Data
		json.Unmarshal(data, &resource)

		addToIndex("authorships", "authorship", resource)
	}
}

func persistResources(dryRun bool, typeName string) {
	if dryRun {
		switch typeName {
		case "people":
			listPeople()
		case "affiliations":
			listPositions()
		case "educations":
			listEducations()
		case "grants":
			listGrants()
		case "publications":
			listPublications()
		}
	} else {
		switch typeName {
		case "people":
			makePeopleIndex() /* won't make if already exists */
			addPeople()
		case "affiliations":
			makeAffiliationsIndex()
			addAffiliations()
		case "educations":
			makeEducationsIndex()
			addEducations()
		case "grants":
			makeGrantsIndex()
			addGrants()
			makeFundingRolesIndex()
			addFundingRoles()
		case "funding-roles":
			makeFundingRolesIndex()
			addFundingRoles()
		case "publications":
			makePublicationsIndex()
			addPublications()
			makeAuthorshipsIndex()
			addAuthorships()
		case "authorships":
			makeAuthorshipsIndex()
			addAuthorships()
		case "all":
			// people
			makePeopleIndex()
			addPeople()
			//affilations
			makeAffiliationsIndex()
			addAffiliations()
			// educations
			makeEducationsIndex()
			addEducations()
			// grants
			makeGrantsIndex()
			addGrants()
			makeFundingRolesIndex()
			addFundingRoles()
			// publications
			makePublicationsIndex()
			addPublications()
			makeAuthorshipsIndex()
			addAuthorships()
		}
	}
}

var conf widgets_import.Config

func main() {
	start := time.Now()
	var err error
	var configFile string
	flag.StringVar(&configFile, "config", "./config.toml", "a config filename")

	dryRun := flag.Bool("dry-run", false, "just examine resources to be saved")
	remove := flag.Bool("remove", false, "remove existing records")
	typeName := flag.String("type", "people", "type of records to import")

	flag.Parse()

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

	// NOTE: elastic client is supposed to be long-lived
	// see https://github.com/olivere/elastic/blob/release-branch.v6/client.go
	client, err = elastic.NewClient(elastic.SetURL(conf.Elastic.Url))
	if err != nil {
		panic(err)
	}

	// NOTE: either remove OR add?
	if *remove {
		clearResources(*typeName)
	} else {
		persistResources(*dryRun, *typeName)
	}

	defer db.Close()
	defer client.Stop()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
