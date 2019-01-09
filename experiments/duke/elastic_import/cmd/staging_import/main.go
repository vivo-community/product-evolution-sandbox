package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/OIT-ads-web/widgets_import"

	"io/ioutil"
	"log"
	//"net/http"
	"os"
	"strings"
	"sync"
	"time"

	//"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	//"github.com/qri-io/jsonschema"
	"github.com/xeipuuv/gojsonschema"
)

var conf widgets_import.Config
var db *sqlx.DB
var psqlInfo string

func GetConnection() *sqlx.DB {
	return db
}

//Created  time.Time

//https://stackoverflow.com/questions/28800672/how-to-add-new-methods-to-an-existing-type-in-go
type Person widgets_import.Person
type Publication widgets_import.Publication
type FundingRole widgets_import.FundingRole
type Grant widgets_import.Grant
type Education widgets_import.Education
type Affiliation widgets_import.Affiliation
type Authorship widgets_import.Authorship

func (person Person) URI() string           { return person.Uri }
func (publication Publication) URI() string { return publication.Uri }
func (role FundingRole) URI() string        { return role.Uri }
func (grant Grant) URI() string             { return grant.Uri }
func (education Education) URI() string     { return education.Uri }
func (affiliation Affiliation) URI() string { return affiliation.Uri }
func (authorship Authorship) URI() string   { return authorship.Uri }

//https://blog.chewxy.com/2018/03/18/golang-interfaces/
type UriAddressable interface {
	URI() string
}

func DeriveUri(u UriAddressable) string { return u.URI() }

func retrieveType(typeName string) []widgets_import.StagingResource {
	db = GetConnection()
	resources := []widgets_import.StagingResource{}

	err := db.Select(&resources, "SELECT id, type, data FROM staging WHERE type =  $1", typeName)
	if err != nil {
		log.Fatalln(err)
	}
	return resources
}

func listType(typeName string) {
	db = GetConnection()
	resources := []widgets_import.StagingResource{}

	err := db.Select(&resources, "SELECT id, type, data FROM staging WHERE type =  $1",
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

//https://stackoverflow.com/questions/2377881/how-to-get-a-md5-hash-from-a-string-in-golang
func makeHash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func resourceExists(uri string, typeName string) bool {
	var exists bool
	db = GetConnection()
	sqlExists := `SELECT EXISTS (SELECT uri FROM RESOURCES where (uri = $1 AND type =$2))`
	db.Get(&exists, sqlExists, uri, typeName)
	return exists
}

// never used?
/*
func addResource(obj interface{}, uri string, typeName string) {
	fmt.Printf(">ADD:%v\n", uri)
	db = GetConnection()

	str, err := json.Marshal(obj)
	if err != nil {
		log.Fatalln(err)
	}
	hash := makeHash(string(str))

	res := &widgets_import.Resource{Uri: uri,
		Type:  typeName,
		Hash:  hash,
		Data:  str,
		DataB: str}

	tx := db.MustBegin()
	sql := `INSERT INTO resources (uri, type, hash, data, data_b) 
	      VALUES (:uri, :type, :hash, :data, :data_b)`
	_, err = tx.NamedExec(sql, res)
	if err != nil {
		log.Fatalln(">ERROR(INSERT):%v", err)
	}
	tx.Commit()
}
*/

// return err ??
func saveResource(obj interface{}, uri string, typeName string) (err error) {
	str, err := json.Marshal(obj)
	if err != nil {
		log.Fatalln(err)
	}

	db = GetConnection()
	hash := makeHash(string(str))

	found := widgets_import.Resource{}
	res := &widgets_import.Resource{Uri: uri,
		Type:  typeName,
		Hash:  hash,
		Data:  str,
		DataB: str}

	findSql := `SELECT uri, type, hash, data, data_b  FROM resources 
	  WHERE (uri = $1 AND type = $2)`

	err = db.Get(&found, findSql, uri, typeName)

	tx := db.MustBegin()
	// error means not found - sql.ErrNoRows
	if err != nil {
		// NOTE: assuming the error means it doesn't exist
		fmt.Printf(">ADD:%v\n", res.Uri)
		sql := `INSERT INTO resources (uri, type, hash, data, data_b) 
	      VALUES (:uri, :type, :hash, :data, :data_b)`
		_, err := tx.NamedExec(sql, res)
		if err != nil {
			log.Fatalln(">ERROR(INSERT):%v", err)
		}
	} else {

		if strings.Compare(hash, found.Hash) == 0 {
			fmt.Printf(">SKIPPING:%v\n", found.Uri)
		} else {
			fmt.Printf(">UPDATE:%v\n", found.Uri)
			sql := `UPDATE resources 
	          set uri = :uri, 
		      type = :type, 
		      hash = :hash, 
		      data = :data, 
		      data_b = :data_b,
		      updated_at = NOW()
		      WHERE uri = :uri and type = :type`
			_, err := tx.NamedExec(sql, res)

			if err != nil {
				log.Fatalln(">ERROR(UPDATE):%v", err)
			}
		}
	}

	tx.Commit()
	return err
}

func markInvalidInStaging(res widgets_import.StagingResource) {
	db = GetConnection()

	//key := PrimaryKey{Id: id, Type: typeName}
	tx := db.MustBegin()
	fmt.Printf(">UPDATE:%v\n", res.Id)
	sql := `UPDATE staging
	    set is_valid = FALSE, 
		WHERE id = :id and type = :type`
	_, err := tx.NamedExec(sql, res)

	if err != nil {
		log.Fatalln(">ERROR(UPDATE):%v", err)
	}
	tx.Commit()
}

func deleteFromStaging(res widgets_import.StagingResource) {
	db = GetConnection()
	sql := `DELETE from staging WHERE id = :id AND type = :type`

	tx := db.MustBegin()
	tx.NamedExec(sql, res)

	log.Println(sql)
	err := tx.Commit()
	if err != nil {
		log.Fatalln(">ERROR(DELETE):%v", err)
	}
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
		// NOTE: organization only come from Positions (now)
		sql += " WHERE type='Position' or type ='Organization'"
	case "grants":
		sql += " WHERE type='Grant' or type='FundingRole'"
	case "publications":
		sql += " WHERE type='Publication' or type='Authorship'"
	case "educations":
		// NOTE: institutions only come from Educations (now)
		sql += " WHERE type='Education' or type='Institution'"
	case "all": // noop
	}
	tx := db.MustBegin()
	tx.MustExec(sql)

	log.Println(sql)
	err := tx.Commit()
	if err != nil {
		log.Fatalln(">ERROR(DELETE):%v", err)
	}
}

func loadSchema(typeName string) *gojsonschema.Schema {
	b, err := ioutil.ReadFile(fmt.Sprintf("schemas/%s.schema.json", typeName)) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	schemaDef := string(b)
	loader1 := gojsonschema.NewStringLoader(schemaDef)
	schema, err := gojsonschema.NewSchema(loader1)

	if err != nil {
		fmt.Println("could not load schema")
		panic(err)
	}
	return schema
}

func validate(schema *gojsonschema.Schema, data string) bool {
	docLoader := gojsonschema.NewStringLoader(data)
	result, err := schema.Validate(docLoader)

	if err != nil {
		fmt.Sprintf("error validating\n")
		return false
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
		if err != nil {
			fmt.Printf("- %s\n", err)
		}
		return true
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, err := range result.Errors() {
			// Err implements the ResultError interface
			fmt.Printf("- %s\n", err)
		}
		return false
	}
}

func addPeople() {
	schema := loadSchema("person")
	people := retrieveType("Person")
	for _, element := range people {
		//resource := widgets_import.Person{}
		resource := Person{}
		data := element.Data
		json.Unmarshal(data, &resource)

		uri := DeriveUri(resource)
		fmt.Println(uri)

		valid := validate(schema, string(data))
		if valid {
			err := saveResource(resource, uri, "Person")
			if err != nil {
				fmt.Printf("- %s\n", err)
			}
		} else {
			//markInvalidInStaging(element)
		}
	}
}

func addAffiliations() {
	schema := loadSchema("affiliation")
	positions := retrieveType("Position")
	for _, element := range positions {
		resource := Affiliation{}
		data := element.Data
		json.Unmarshal(data, &resource)

		uri := DeriveUri(resource)
		fmt.Println(uri)
		valid := validate(schema, string(data))
		if valid {
			err := saveResource(resource, uri, "Affiliation")
			if err != nil {
				fmt.Printf("- %s\n", err)
			}
		}
	}
}

func addEducations() {
	schema := loadSchema("education")
	educations := retrieveType("Education")
	for _, element := range educations {
		resource := Education{}
		data := element.Data
		json.Unmarshal(data, &resource)

		uri := DeriveUri(resource)
		fmt.Println(uri)
		valid := validate(schema, string(data))
		if valid {
			err := saveResource(resource, uri, "Education")
			if err != nil {
				fmt.Printf("- %s\n", err)
			}
		}
	}
}

func addGrants() {
	schema := loadSchema("grant")
	grants := retrieveType("Grant")
	for _, element := range grants {
		resource := Grant{}
		data := element.Data
		json.Unmarshal(data, &resource)

		uri := DeriveUri(resource)
		fmt.Println(uri)
		valid := validate(schema, string(data))
		if valid {
			err := saveResource(resource, uri, "Grant")
			if err != nil {
				fmt.Printf("- %s\n", err)
			}

		}
	}
}

func addFundingRoles() {
	schema := loadSchema("funding-role")
	fundingRoles := retrieveType("FundingRole")
	for _, element := range fundingRoles {
		resource := FundingRole{}
		data := element.Data
		json.Unmarshal(data, &resource)

		uri := DeriveUri(resource)
		fmt.Println(uri)
		valid := validate(schema, string(data))
		if valid {
			err := saveResource(resource, uri, "FundingRole")
			if err != nil {
				fmt.Printf("- %s\n", err)
			}
		}
	}
}

func addPublications() {
	schema := loadSchema("publication")
	publications := retrieveType("Publication")
	for _, element := range publications {
		//resource := widgets_import.Publication{}
		resource := Publication{}

		data := element.Data
		json.Unmarshal(data, &resource)

		uri := DeriveUri(resource)
		fmt.Println(uri)
		valid := validate(schema, string(data))
		if valid {
			err := saveResource(resource, uri, "Publication")
			if err != nil {
				fmt.Printf("- %s\n", err)
			}
		}
	}
}

func addAuthorships() {
	schema := loadSchema("authorship")
	authorships := retrieveType("Authorship")
	for _, element := range authorships {
		resource := Authorship{}
		data := element.Data
		json.Unmarshal(data, &resource)

		uri := DeriveUri(resource)
		fmt.Println(uri)
		valid := validate(schema, string(data))
		if valid {
			err := saveResource(resource, uri, "Authorship")
			if err != nil {
				fmt.Printf("- %s\n", err)
			}
		}
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
		case "all":
			listPeople()
			listPositions()
			listEducations()
			listGrants()
			listPublications()
		}
	} else {
		switch typeName {
		case "people":
			addPeople()
		case "affiliations":
			addAffiliations()
		case "educations":
			addEducations()
		case "grants":
			addGrants()
			addFundingRoles()
		case "funding-roles":
			addFundingRoles()
		case "publications":
			// parallelize?
			addPublications()
			addAuthorships()
		case "authorships":
			addAuthorships()
		case "all":
			//25.860895353s
			/*
				addPeople()
				addAffiliations()
				addEducations()
				addGrants()
				addFundingRoles()
				addPublications()
				addAuthorships()
			*/
			// trying to let it do things
			// in goroutines
			wg.Add(7)
			/// this doesn't stop itself
			/*
							defer wg.Done()

							go addPeople()
							go addAffiliations()
							go addEducations()
				            go addGrants()
							go addFundingRoles()
				            go addPublications()
							go addAuthorships()
			*/
			//17.970656632s
			// 1.people
			go func() {
				defer wg.Done()
				addPeople()
			}()
			// 2. affilations
			go func() {
				defer wg.Done()
				addAffiliations()
			}()
			// 3. educations
			go func() {
				defer wg.Done()
				addEducations()
			}()
			// 4. grants
			go func() {
				defer wg.Done()
				addGrants()
			}()
			// 5. funding-roles
			go func() {
				defer wg.Done()
				addFundingRoles()
			}()
			// 6. publications
			go func() {
				defer wg.Done()
				addPublications()
			}()
			// 7. authorships
			go func() {
				defer wg.Done()
				addAuthorships()
			}()

			wg.Wait()
		}
	}
}

var wg sync.WaitGroup

// import from staging table -> resources table
// go through jsonschema validate
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

	if !resourceTableExists() {
		makeResourceSchema()
	}

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
