package psql

import (
	"github.com/OIT-ads-web/widgets_import"
	//"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func RetrieveType(typeName string, updates bool) []widgets_import.Resource {
	db := GetConnection()
	resources := []widgets_import.Resource{}

	var err error
	if updates {
		// TODO: ideally would need to record time last run somewhere
		yesterday := time.Now().AddDate(0, 0, -1)
        rounded := time.Date(yesterday.Year(), yesterday.Month(), 
		    yesterday.Day(), 0, 0, 0, 0, yesterday.Location())

		sql := `SELECT uri, type, hash, data 
		FROM resources 
		WHERE type =  $1 and updated_at >= $2
      `
		err = db.Select(&resources, sql, typeName, rounded)
	} else {
		err = db.Select(&resources, "SELECT uri, type, hash, data FROM resources WHERE type =  $1", typeName)
	}

	if err != nil {
		log.Fatalln(err)
	}
	return resources
}

func ListType(typeName string, updates bool) {
	db := GetConnection()
	resources := []widgets_import.Resource{}

	var err error
	if updates {
		// TODO: ideally would need to record time last run somewhere
		yesterday := time.Now().AddDate(0, 0, -1)
		//yesterday := time.Now()

        rounded := time.Date(yesterday.Year(), yesterday.Month(), 
		    yesterday.Day(), 0, 0, 0, 0, yesterday.Location())

		sql := `SELECT uri, type, hash, data 
		FROM resources 
		WHERE type =  $1 and updated_at >= $2
      `
		err = db.Select(&resources, sql, typeName, rounded)
	} else {
		err = db.Select(&resources, "SELECT uri, type, hash, data FROM resources WHERE type =  $1", typeName)
	}

	//err := db.Select(&resources, "SELECT uri, type, hash, data FROM resources WHERE type =  $1",
	//	typeName)
	for _, element := range resources {
		log.Println(element)
		// element is the element from someSlice for where we are
	}
	log.Printf("******* count = %d ********\n", len(resources))

	if err != nil {
		log.Fatalln(err)
	}
}

func ListPeople(updated bool) {
	ListType("Person", updated)
}

func ListPositions(updated bool) {
	ListType("Position", updated)
}

func ListEducations(updated bool) {
	ListType("Education", updated)
}

func ListGrants(updated bool) {
	ListType("Grant", updated)
}

func ListFundingRoles(updated bool) {
	ListType("FundingRole", updated)
}

func ListPublications(updated bool) {
	ListType("Publication", updated)
}

// just a wrapper
func ListAffiliations(updated bool) {
	ListPositions(updated)
}
