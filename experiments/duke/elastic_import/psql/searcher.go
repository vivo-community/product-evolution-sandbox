package psql

import (
	"github.com/OIT-ads-web/widgets_import"
	//"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func RetrieveType(typeName string) []widgets_import.Resource {
	db := GetConnection()
	resources := []widgets_import.Resource{}

	err := db.Select(&resources, "SELECT uri, type, hash, data FROM resources WHERE type =  $1", typeName)
	if err != nil {
		log.Fatalln(err)
	}
	return resources
}

func ListType(typeName string) {
	db := GetConnection()
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

func ListPeople() {
	ListType("Person")
}

func ListPositions() {
	ListType("Position")
}

func ListEducations() {
	ListType("Education")
}

func ListGrants() {
	ListType("Grant")
}

func ListFundingRoles() {
	ListType("FundingRole")
}

func ListPublications() {
	ListType("Publication")
}

// just a wrapper
func ListAffiliations() {
	ListType("Position")
}

