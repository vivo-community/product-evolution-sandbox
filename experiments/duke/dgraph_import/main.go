package main

import (
	"context"
	"encoding/json"
	"os"

	//"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

type database struct {
	Server   string
	Port     int
	Database string
	User     string
	Password string
}

type Config struct {
	Database database `toml:"database"`
	DGraph   dgraph   `toml:"dgraph"`
}

type dgraph struct {
	Url string
}

// ********** database json column structs:
// NOTE: this is *not* an independent resource, should it be?
type Keyword struct {
	// not sure an 'id' makes sense - they are like #mesh, LOC etc...
	Uri   string
	Label string
}

// neither is this -in RDF it has to be, but seems like overkill
type DateResolution struct {
	DateTime   string
	Resolution string
}

type ResourceFundingRole struct {
	Id       string
	Uri      string
	GrantId  string
	PersonId string
	RoleName string
}

type ResourceGrant struct {
	Id                      string
	Uri                     string
	Label                   string
	PrincipalInvestigatorId string
	Start                   DateResolution
	End                     DateResolution
}

type ResourcePerson struct {
	Id                string
	Uri               string
	AlternateId       string
	FirstName         string
	LastName          string
	MiddleName        *string
	PrimaryTitle      string
	ImageUri          string
	ImageThumbnailUri string
	Type              string
	Overview          string
	Keywords          []Keyword
}

type PersonKeyword struct {
	// not sure an 'id' makes sense - they are like #mesh, LOC etc...
	Uri   string `json:"keywordUri,omitempty"`
	Label string `json:"keywordLabel,omitempty"`
}

type Person struct {
	Id                string          `json:"personId,omitempty"`
	Uri               string          `json:"personUri,omitempty"`
	AlternateId       string          `json:"personAlternateId,omitempty"`
	FirstName         string          `json:"personFirstName,omitempty"`
	LastName          string          `json:"personLastName,omitempty"`
	MiddleName        *string         `json:"personMiddleName,omitempty"`
	PrimaryTitle      string          `json:"personPrimaryTitle,omitempty"`
	ImageUri          string          `json:"personImageUri,omitempty"`
	ImageThumbnailUri string          `json:"personImageThumbnailUri,omitempty"`
	Type              string          `json:"personType,omitempty"`
	Overview          string          `json:"personOverview,omitempty"`
	Keywords          []PersonKeyword `json:"personKeywords,omitempty"`
}
type ResourcePosition struct {
	Id                string
	Uri               string
	PersonId          string
	Label             string
	Start             DateResolution
	OrganizationId    string
	OrganizationLabel string
}

type ResourceInstitution struct {
	Id    string
	Uri   string
	Label string
}

type ResourceEducation struct {
	Id               string
	Uri              string
	PersonId         string
	Label            string
	InsitutionId     string
	InstitutionLabel string
}

type ResourceAuthorship struct {
	Id             string
	Uri            string
	PublicationId  string
	PersonId       string
	AuthorshipType string
}

type ResourcePublication struct {
	Id                  string
	Uri                 string
	Label               string
	AuthorList          string
	Doi                 string
	PublishedIn         string
	PublicationVenueUri string
}

type ResourceOrganization struct {
	Id    string
	Uri   string
	Label string
}

/*** end database json column object maps */

// this is the raw structure in the database
// two json columms:
// * 'data' can be used for change comparison with hash
// * 'data_b' can be used for searches
type Resource struct {
	Uri   string         `db:"uri"`
	Type  string         `db:"type"`
	Hash  string         `db:"hash"`
	Data  types.JSONText `db:"data"`
	DataB types.JSONText `db:"data_b"`
}

// ********** end database json structs
func getList(typeName string) []Resource {
	db = GetConnection()
	resources := []Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1",
		typeName)

	if err != nil {
		panic(err)
	}
	return resources
}

func listType(typeName string) {
	db = GetConnection()
	resources := []Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data, data_b FROM resources WHERE type =  $1",
		typeName)
	for _, element := range resources {
		log.Println(element)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func makePeopleIndex() {
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	op := &api.Operation{}
	op.Schema = `
        personId: string @index(exact) .
		personFirstName: string @index(exact) .
		personLastName: string @index(exact) .
	    personAlternateId: string @index(exact) .
	    personPrimaryTitle: string @index(exact) .
	    personOverview: string @index(fulltext) .
	`

	ctx := context.Background()
	err = dg.Alter(ctx, op)
	if err != nil {
		log.Fatal(err)
	}
}

func addPeople() {
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	ctx := context.Background()

	mu := &api.Mutation{
		CommitNow: true,
	}

	people := getList("Person")

	for _, row := range people {
		log.Println(row)
		personJson := row.DataB

		// 1. get db stract
		var resource ResourcePerson
		json.Unmarshal(personJson, &resource)

		// 2. copy PersonResource->Person (for dgraph)
		// conversion?
		var personKeywords []PersonKeyword
		for _, keyword := range resource.Keywords {
			// implicit conversion (since they are exactly the same)
			personKeywords = append(personKeywords, PersonKeyword(keyword))
		}

		person := Person{resource.Id,
			resource.Uri,
			resource.AlternateId,
			resource.FirstName,
			resource.LastName,
			resource.MiddleName,
			resource.PrimaryTitle,
			resource.ImageUri,
			resource.ImageThumbnailUri,
			resource.Type,
			resource.Overview,
			personKeywords}

		// then back to Json
		json, err := json.Marshal(person)
		if err != nil {
			continue
		}
		log.Println(json)
		mu.SetJson = json
		_, err = dg.NewTxn().Mutate(ctx, mu)
		if err != nil {
			log.Fatal(err)
		}
		//variables := map[string]string{"$id": assigned.Uids["blank-0"]}
		//log.Println(variables)
	}
}

func clearResources(typeName string) {
	switch typeName {
	case "people":
		fmt.Println("not implemented")
	}
}

func listPeople() {
	listType("Person")
}

func persistResources(dryRun bool, typeName string) {
	if dryRun {
		switch typeName {
		case "people":
			listPeople()
		}
	} else {
		switch typeName {
		case "people":
			makePeopleIndex()
			addPeople()
		}
	}
}

var psqlInfo string
var db *sqlx.DB
var conf Config

func GetConnection() *sqlx.DB {
	return db
}

func main() {
	start := time.Now()
	var err error
	var configFile string
	flag.StringVar(&configFile, "config", "./config.toml", "a config filename")
	typeName := flag.String("type", "people", "type of records to import")
	dryRun := flag.Bool("dry-run", false, "just examine resources to be saved")
	remove := flag.Bool("remove", false, "remove existing records")

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

	if *remove {
		clearResources(*typeName)
	} else {
		persistResources(*dryRun, *typeName)
	}

	defer db.Close()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
