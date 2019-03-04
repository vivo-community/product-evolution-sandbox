package elastic

import (
	"fmt"
	"log"
	"github.com/olivere/elastic"
	"context"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/OIT-ads-web/widgets_import"
	//"github.com/OIT-ads-web/widgets_import/templates"
)

func addToIndex(index string, typeName string, id string, obj interface{}) {
	ctx := context.Background()
	client := GetClient()

	get1, err := client.Get().
		Index(index).
		Type(typeName).
		Id(id).
		Do(ctx)

	switch {
	case elastic.IsNotFound(err):
		put1, err := client.Index().
			Index(index).
			Type(typeName).
			Id(id).
			BodyJson(obj).
			Do(ctx)

		if err != nil {
			panic(err)
		}

		fmt.Printf("ADDED %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		spew.Println(obj)
		return
	case elastic.IsConnErr(err):
		panic(err)
	case elastic.IsTimeout(err):
		panic(err)
	case err != nil:
		panic(err)
	}

	if get1.Found {
		update1, err := client.Update().
			Index(index).
			Type(typeName).
			Id(id).
			Doc(obj).
			Do(ctx)

		if err != nil {
			panic(err)
		}

		fmt.Printf("UPDATED %s to index %s, type %s\n", update1.Id, update1.Index, update1.Type)
	}

	if err != nil {
		// Handle error
		panic(err)
	}
	spew.Println(obj)
}

func partialUpdate(index string, typeName string, id string, prop string, obj interface{}) {
	ctx := context.Background()
	client := GetClient()

	get1, err := client.Get().
		Index(index).
		Type(typeName).
		Id(id).
		Do(ctx)

	switch {
	case elastic.IsNotFound(err):
		// NOTE: in theory we could add without source doc
		fmt.Printf("no doc id=%s found to append to\n", id)
		//fmt.Printf("ADDED %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		return
	case elastic.IsConnErr(err):
		panic(err)
	case elastic.IsTimeout(err):
		panic(err)
	case err != nil:
		panic(err)
	}

	if get1.Found {
		update1, err := client.Update().
			Index(index).
			Type(typeName).
			Id(id).
			//Doc(obj).
			// replace all of prop ??...
			Doc(map[string]interface{}{prop: obj}).
			DetectNoop(true).
			Do(ctx)

		if err != nil {
			panic(err)
		}

		fmt.Printf("UPDATED %s to index %s, type %s\n", update1.Id, update1.Index, update1.Type)
	}

	if err != nil {
		// Handle error
		panic(err)
	}
	spew.Println(obj)
}

func clearIndex(name string) {
	ctx := context.Background()

	client := GetClient()

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

func ClearPeopleIndex() {
	clearIndex("people")
}

func ClearAffiliationsIndex() {
	clearIndex("affiliations")
}

func ClearEducationsIndex() {
	clearIndex("educations")
}

func ClearGrantsIndex() {
	clearIndex("grants")
}

func ClearFundingRolesIndex() {
	clearIndex("funding-roles")
}

func ClearPublicationsIndex() {
	clearIndex("publications")
}

func ClearAuthorshipsIndex() {
	clearIndex("authorships")
}

// NOTE: 'mappingJson' is just a json string plugged into template
func makeIndex(name string, mappingJson string) {
	ctx := context.Background()
	
	client := GetClient()

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(name).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(name).BodyString(mappingJson).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
}

func MakePeopleIndex(mapping string) {
	makeIndex("people", mapping)
}

func MakeGrantsIndex(mapping string) {
	makeIndex("grants", mapping)
}

func MakeFundingRolesIndex(mapping string) {
	makeIndex("funding-roles", mapping)
}

func MakePublicationsIndex(mapping string) {
	makeIndex("publications", mapping)
}

func MakeAuthorshipsIndex(mapping string) {
	makeIndex("authorships", mapping)
}

func AddPeople(people []widgets_import.Resource) {
	for _, element := range people {
		resource := widgets_import.Person{}
		data := element.Data
		json.Unmarshal(data, &resource)

		addToIndex("people", "person", resource.Id, resource)
	}
}

func AddAffiliationsToPeople(positions []widgets_import.Resource) {
	// need to group by personId
	collections := make(map[string][]widgets_import.Affiliation)

	for _, element := range positions {
		resource := widgets_import.Affiliation{}
		data := element.Data
		json.Unmarshal(data, &resource)

		collections[resource.PersonId] = append(collections[resource.PersonId], resource)
	}

	for key, value := range collections {
		partialUpdate("people", "person", key, "affiliationList", value)
	}
}

func AddEducationsToPeople(educations []widgets_import.Resource) {
	// need to group by personId
	collections := make(map[string][]widgets_import.Education)

	for _, element := range educations {
		resource := widgets_import.Education{}
		data := element.Data
		json.Unmarshal(data, &resource)

		collections[resource.PersonId] = append(collections[resource.PersonId], resource)
	}
	for key, value := range collections {
		partialUpdate("people", "person", key, "educationList", value)
	}
}

func AddGrants(grants []widgets_import.Resource) {
	for _, element := range grants {
		resource := widgets_import.Grant{}
		data := element.Data
		json.Unmarshal(data, &resource)

		addToIndex("grants", "grant", resource.Id, resource)
	}
}

func AddFundingRoles(fundingRoles []widgets_import.Resource) {
	for _, element := range fundingRoles {
		resource := widgets_import.FundingRole{}
		data := element.Data
		json.Unmarshal(data, &resource)

		addToIndex("funding-roles", "funding-role", resource.Id, resource)
	}
}

func AddPublications(publications []widgets_import.Resource) {
	for _, element := range publications {
		resource := widgets_import.Publication{}
		data := element.Data
		json.Unmarshal(data, &resource)

		addToIndex("publications", "publication", resource.Id, resource)
	}
}

func AddAuthorships(authorships []widgets_import.Resource) {
	for _, element := range authorships {
		resource := widgets_import.Authorship{}
		data := element.Data
		json.Unmarshal(data, &resource)

		addToIndex("authorships", "authorship", resource.Id, resource)
	}
}

