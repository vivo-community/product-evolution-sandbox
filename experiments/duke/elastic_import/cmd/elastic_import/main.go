package main

import (
	"flag"
	"fmt"
	wi "github.com/OIT-ads-web/widgets_import"
	"github.com/OIT-ads-web/widgets_import/elastic"
	"github.com/OIT-ads-web/widgets_import/psql"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
	"sync"
	"time"
)

var conf wi.Config
var wg sync.WaitGroup

func preview(typeName string, updated bool) {
	switch typeName {
	case "people":
		psql.ListPeople(updated)
		mapping, err := elastic.PersonMapping()
		if err != nil {
		    fmt.Printf("error %s\n", err)
			break
		}
		wi.Preview(mapping)
	case "publications":
		psql.ListPublications(updated)
		mapping, err := elastic.PublicationMapping()
		if err != nil {
		    fmt.Printf("error %s\n", err)
			break
		}
		wi.Preview(mapping)
	case "grants":
		psql.ListGrants(updated)
		mapping, err := elastic.GrantMapping()
		if err != nil {
		    fmt.Printf("error %s\n", err)
			break
		}
		wi.Preview(mapping)
	case "all":
		fmt.Println("no option to preview 'all'")
	}
}

func clearIndexes(typeName string) {
	switch typeName {
	case "people":
		elastic.ClearPeopleIndex()
	case "grants":
		elastic.ClearGrantsIndex()
		elastic.ClearFundingRolesIndex()
	case "publications":
		elastic.ClearPublicationsIndex()
		elastic.ClearAuthorshipsIndex()
	case "all":
		elastic.ClearPeopleIndex()
		elastic.ClearGrantsIndex()
		elastic.ClearFundingRolesIndex()
		elastic.ClearPublicationsIndex()
		elastic.ClearAuthorshipsIndex()
	}
}

func persistResources(typeName string, updates bool) {
	fmt.Println("only updates = %t\n", updates)
	switch typeName {
	case "people":
		mapping, err := elastic.PersonMapping()
		if err != nil {
			fmt.Println(err)
			break
		}
		elastic.MakePeopleIndex(mapping)
		people := psql.RetrieveType("Person", updates)
		elastic.AddPeople(people)
	case "affiliations":
		affiliations := psql.RetrieveType("Affiliation", updates)
		elastic.AddAffiliationsToPeople(affiliations)
	case "educations":
		educations := psql.RetrieveType("Education", updates)
		elastic.AddEducationsToPeople(educations)
	case "grants":
		mapping, err := elastic.GrantMapping()
		if err != nil {
			fmt.Println(err)
			break
		}

		elastic.MakeGrantsIndex(mapping)
		grants := psql.RetrieveType("Grant", updates)
		elastic.AddGrants(grants)

		mapping, err = elastic.FundingRoleMapping()
		if err != nil {
			fmt.Println(err)
			break
		}

		elastic.MakeFundingRolesIndex(mapping)
		roles := psql.RetrieveType("FundingRole", updates)
		elastic.AddFundingRoles(roles)
	case "funding-roles":
		mapping, err := elastic.FundingRoleMapping()
		if err != nil {
			fmt.Println(err)
			break
		}

		elastic.MakeFundingRolesIndex(mapping)

		roles := psql.RetrieveType("FundingRole", updates)
		elastic.AddFundingRoles(roles)
	case "publications":
		mapping, err := elastic.PublicationMapping()
		if err != nil {
			fmt.Println(err)
			break
		}

		elastic.MakePublicationsIndex(mapping)
		publications := psql.RetrieveType("Publication", updates)
		elastic.AddPublications(publications)

		mapping, err = elastic.AuthorshipMapping()
		if err != nil {
			fmt.Println(err)
			break
		}

		elastic.MakeAuthorshipsIndex(mapping)
		authorships := psql.RetrieveType("Authorship", updates)
		elastic.AddAuthorships(authorships)
	case "authorships":
		mapping, err := elastic.GrantMapping()
		if err != nil {
			fmt.Println(err)
			break
		}

		elastic.MakeAuthorshipsIndex(mapping)

		authorships := psql.RetrieveType("Authorship", updates)
		elastic.AddAuthorships(authorships)
	case "all":
		mapping, err := elastic.PersonMapping()
		if err != nil {
			fmt.Println(err)
			break
		}
		elastic.MakePeopleIndex(mapping)

		mapping, err = elastic.PublicationMapping()
		if err != nil {
			fmt.Println(err)
			break
		}
		elastic.MakePublicationsIndex(mapping)

		mapping, err = elastic.AuthorshipMapping()
		if err != nil {
			fmt.Println(err)
			break
		}
		elastic.MakeAuthorshipsIndex(mapping)

		mapping, err = elastic.GrantMapping()
		if err != nil {
			fmt.Println(err)
			break
		}
		elastic.MakeGrantsIndex(mapping)

		mapping, err = elastic.FundingRoleMapping()
		if err != nil {
			fmt.Println(err)
			break
		}
		elastic.MakeFundingRolesIndex(mapping)

		wg.Add(7)
		// 1.people
		go func() {
			defer wg.Done()
			people := psql.RetrieveType("Person", updates)
			elastic.AddPeople(people)
		}()
		// 2. affilations
		go func() {
			defer wg.Done()
			affiliations := psql.RetrieveType("Affiliation", updates)
			elastic.AddAffiliationsToPeople(affiliations)
		}()
		// 3. educations
		go func() {
			defer wg.Done()
			educations := psql.RetrieveType("Education", updates)
			elastic.AddEducationsToPeople(educations)
		}()
		// 4. grants
		go func() {
			defer wg.Done()
			grants := psql.RetrieveType("Grant", updates)
			elastic.AddGrants(grants)
		}()
		// 5. funding-roles
		go func() {
			defer wg.Done()
			roles := psql.RetrieveType("FundingRole", updates)
			elastic.AddFundingRoles(roles)
		}()
		// 6. publications
		go func() {
			defer wg.Done()
			publications := psql.RetrieveType("Publication", updates)
			elastic.AddPublications(publications)
		}()
		// 7. authorships
		go func() {
			defer wg.Done()
			authorships := psql.RetrieveType("Authorship", updates)
			elastic.AddAuthorships(authorships)
		}()

		wg.Wait()
	}
}

func main() {
	start := time.Now()

	if os.Getenv("ENVIRONMENT") == "development" {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")

		value, exists := os.LookupEnv("CONFIG_PATH")
		if exists {
			viper.AddConfigPath(value)
		}

		viper.ReadInConfig()
	} else {
		replacer := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(replacer)
		viper.BindEnv("database.server")
		viper.BindEnv("database.port")
		viper.BindEnv("database.database")
		viper.BindEnv("database.user")
		viper.BindEnv("database.password")
		viper.BindEnv("elastic.url")
		viper.BindEnv("templates.layout")
		viper.BindEnv("templates.include")
	}

	dryRun := flag.Bool("dry-run", false, "just examine resources to be saved")
	remove := flag.Bool("remove", false, "remove existing records")
	typeName := flag.String("type", "people", "type of records to import")
	updates := flag.Bool("updates", true, "only import updated records")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if err := viper.Unmarshal(&conf); err != nil {
		fmt.Printf("could not establish read into conf structure %s\n", err)
		os.Exit(1)
	}

	if err := elastic.MakeClient(conf.Elastic.Url); err != nil {
		fmt.Printf("could not establish elastic client %s\n", err)
		os.Exit(1)
	}

	if err := psql.MakeConnection(conf); err != nil {
		fmt.Printf("could not establish postgresql connection %s\n", err)
		os.Exit(1)
	}

	wi.LoadTemplates(conf)

	// NOTE: elastic client is supposed to be long-lived
	// see https://github.com/olivere/elastic/blob/release-branch.v6/client.go
	//client, err = elastic.NewClient(elastic.SetURL(conf.Elastic.Url))
	//if err != nil {
	//	panic(err)
	//}

	// NOTE: either remove OR add?
	if *remove {
		clearIndexes(*typeName)
	} else {
		if *dryRun {
			preview(*typeName, *updates)
		} else {
			// only updates?
			persistResources(*typeName, *updates)
		}
		// if dryRun -> listResources
		// -> list Mappings ->
		//  else elastic.persistResources( -- get from psql -- )
		//elastic.persistResources(*typeName)
	}

	defer psql.Database.Close()
	defer elastic.Client.Stop()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
