package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/OIT-ads-web/widgets_import"
	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/olivere/elastic"
	"log"
	"os"
	"text/template"
	"time"
)

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
	Label string       `json:"overview"`
	Type  OverviewType `json:"type"`
}

type Person struct {
	Id           string           `json:"id"`
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

type Organization struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

type Institution struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

// NOTE: model doesn't have org as sub-object
type Affiliation struct {
	Id                string `json:"id"`
	Uri               string `json:"uri"`
	PersonId          string `json:"personId"`
	Label             string `json:"label"`
	StartDate         Date   `json:"startDate"`
	OrganizationId    string `json:"organizationId"`
	OrganizationLabel string `json:"organizationLabel"`
}

type Education struct {
	Id          string      `json:"id"`
	Uri         string      `json:"id"`
	Label       string      `json:"label"`
	PersonId    string      `json:"personId"`
	Institution Institution `json:"org" elastic:"type:object"`
}

type FundingRole struct {
	Id       string `json:"id"`
	Uri      string `json:"uri"`
	GrantId  string `json:"grantId"`
	PersonId string `json:"personId"`
	Label    string `json:"label"`
}

type Grant struct {
	Id        string `json:"id"`
	Uri       string `json:"uri"`
	Label     string `json:"label"`
	StartDate Date   `json:"startDate"`
	EndDate   Date   `json:"endDate"`
}

type Authorship struct {
	Id            string `json:"id"`
	Uri           string `json:"uri"`
	PersonId      string `json:"personId"`
	PublicationId string `json:"publicationId"`
	Label         string `json:"label"`
}

type Publication struct {
	Id         string `json:"id"`
	Uri        string `json:"uri"`
	Label      string `json:"label"`
	// NOTE: this is supposed to be an array
	AuthorList string `json:"authorList"`
	Doi        string `json:"doi"`
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

// TODO: publications not quite worked out
const publicationMapping = `
"publication":{
	"properties":{
		"id":        { "type": "text" },
		"uri":       { "type": "text" },
		"label":     { "type": "text" }
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

func listType(typeName string) {
	db = GetConnection()
	resources := []widgets_import.Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1",
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
		// Handle error
		panic(err)
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
		//clearFundingRolesIndex()
	case "publications":
		clearPublicationsIndex()
	case "all":
		clearPeopleIndex()
		clearAffiliationsIndex()
		clearEducationsIndex()
		clearPublicationsIndex()
		clearGrantsIndex()
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
	db = GetConnection()
	resources := []widgets_import.Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1", "Person")
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
		person := Person{resource.Id, resource.Uri, resource.AlternateId, resource.PrimaryTitle,
			name, image, personType, overviewList, keywordList}

		addToIndex("people", "person", person)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func makePositionDate(position widgets_import.ResourcePosition) Date {
	return Date{position.Start.DateTime, position.Start.Resolution}
}

func addAffiliations() {
	db = GetConnection()
	resources := []widgets_import.Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1", "Position")
	for _, element := range resources {
		resource := widgets_import.ResourcePosition{}
		data := element.Data
		json.Unmarshal(data, &resource)

		date := makePositionDate(resource)

		affiliation := Affiliation{resource.Id,
			resource.Uri,
			resource.PersonId,
			resource.Label,
			date,
			resource.OrganizationId,
			resource.OrganizationLabel}
		addToIndex("affiliations", "affiliation", affiliation)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func addEducations() {
	db = GetConnection()
	resources := []widgets_import.Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1", "Education")
	for _, element := range resources {
		resource := widgets_import.ResourceEducation{}
		data := element.Data
		json.Unmarshal(data, &resource)

		institution := Institution{resource.InsitutionId, resource.InstitutionLabel}

		education := Education{resource.Id,
			resource.Uri,
			resource.Label,
			resource.PersonId,
			institution}
		addToIndex("educations", "education", education)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func makeGrantDates(grant widgets_import.ResourceGrant) (Date, Date) {
	start := Date{grant.Start.DateTime, grant.Start.Resolution}
	end := Date{grant.End.DateTime, grant.End.Resolution}
	return start, end
}

/*
type ResourceGrant struct {
	Id                      string
	Uri                     string
	Label                   string
	PrincipalInvestigatorId strin
	Start                   DateResolution
	End                     DateResolution
}

into this:

type Grant struct {
	Id        string `json:"id"`
	Uri       string `json:"uri"`
	Label     string `json:"label"`
	StartDate Date   `json:"startDate"`
	EndDate   Date   `json:"startDate"`
}
*/
func addGrants() {
 	db = GetConnection()
	resources := []widgets_import.Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1", "Grant")
	for _, element := range resources {
		resource := widgets_import.ResourceGrant{}
		data := element.Data
		json.Unmarshal(data, &resource)
        start, end := makeGrantDates(resource)

		grant := Grant{resource.Id,
			resource.Uri,
			resource.Label,
			start,
			end}
		addToIndex("grants", "grant", grant)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

/*
type ResourceFundingRole struct {
	Id       string
	Uri      string
	GrantId  string
	PersonId string
	RoleName string
}
into this:

type FundingRole struct {
	Id       string `json:"id"`
	Uri      string `json:"uri"`
	GrantId  string `json:"grantId"`
	PersonId string `json:"personId"`
	Label    string `json:"label"`
}
*/
func addFundingRoles() {
 	db = GetConnection()
	resources := []widgets_import.Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1", "FundingRole")
	for _, element := range resources {
		resource := widgets_import.ResourceFundingRole{}
		data := element.Data
		json.Unmarshal(data, &resource)

		role := FundingRole{resource.Id,
			resource.Uri,
			resource.GrantId,
			resource.PersonId,
			resource.RoleName}
		addToIndex("funding-roles", "funding-role", role)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

/*
type ResourcePublication struct {
	Id         string
	Uri        string
	Label      string
	AuthorList string
	Doi        string
}
*/
func addPublications() {

}
/*
type ResourceAuthorship struct {
	Id            string
	Uri           string
	PublicationId string
	PersonId      string
}
*/
func addAuthorships() {

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

	// NOTE: elastic client is supposed to be long-lived
	// see https://github.com/olivere/elastic/blob/release-branch.v6/client.go
	client, err = elastic.NewClient(elastic.SetURL(conf.Elastic.Url))
	if err != nil {
		panic(err)
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
	defer client.Stop()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
