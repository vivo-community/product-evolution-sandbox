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

func preview(typeName string) {
	switch typeName {
	case "people":
		//psql.ListPeople()
		wi.Preview(elastic.PersonMapping())
	case "affiliations":
		//psql.ListAffiliations()
		wi.Preview(elastic.AffiliationMapping())
	case "publications":
		//psql.ListPublications()
		wi.Preview(elastic.PublicationMapping())
	case "grants":
		//psql.ListGrants()
		wi.Preview(elastic.GrantMapping())

	case "all":
		psql.ListPeople()
		psql.ListAffiliations()
	}
}

func clearIndexes(typeName string) {
	switch typeName {
	case "people":
		elastic.ClearPeopleIndex()
	//case "affiliations":
	//	elastic.ClearAffiliationsIndex()
	//case "educations":
	//	elastic.ClearEducationsIndex()
	case "grants":
		elastic.ClearGrantsIndex()
		elastic.ClearFundingRolesIndex()
	case "publications":
		elastic.ClearPublicationsIndex()
		elastic.ClearAuthorshipsIndex()
	case "all":
		elastic.ClearPeopleIndex()
		//elastic.ClearAffiliationsIndex()
		//elastic.ClearEducationsIndex()
		elastic.ClearGrantsIndex()
		elastic.ClearFundingRolesIndex()
		elastic.ClearPublicationsIndex()
		elastic.ClearAuthorshipsIndex()
	}
}

func persistResources(typeName string) {
	switch typeName {
	case "people":
		elastic.MakePeopleIndex(elastic.PersonMapping())

		people := psql.RetrieveType("Person")
		elastic.AddPeople(people)
	case "affiliations":
		affiliations := psql.RetrieveType("Affiliation")
		elastic.AddAffiliationsToPeople(affiliations)
	case "educations":
		educations := psql.RetrieveType("Education")
		elastic.AddEducationsToPeople(educations)
	case "grants":
		elastic.MakeGrantsIndex(elastic.GrantMapping())
		grants := psql.RetrieveType("Grant")
		elastic.AddGrants(grants)

		elastic.MakeFundingRolesIndex(elastic.FundingRoleMapping())
		roles := psql.RetrieveType("FundingRole")
		elastic.AddFundingRoles(roles)
	case "funding-roles":
		elastic.MakeFundingRolesIndex(elastic.FundingRoleMapping())

		roles := psql.RetrieveType("FundingRole")
		elastic.AddFundingRoles(roles)
	case "publications":
		elastic.MakePublicationsIndex(elastic.PublicationMapping())
		publications := psql.RetrieveType("Publication")
		elastic.AddPublications(publications)

		elastic.MakeAuthorshipsIndex(elastic.AuthorshipMapping())
		authorships := psql.RetrieveType("FundingRole")
		elastic.AddAuthorships(authorships)
	case "authorships":
		elastic.MakeAuthorshipsIndex(elastic.AuthorshipMapping())

		authorships := psql.RetrieveType("Authorship")
		elastic.AddAuthorships(authorships)
	case "all":

		elastic.MakePeopleIndex(elastic.PersonMapping())
		elastic.MakeGrantsIndex(elastic.GrantMapping())
		elastic.MakeFundingRolesIndex(elastic.FundingRoleMapping())
		elastic.MakePublicationsIndex(elastic.PublicationMapping())
		elastic.MakeAuthorshipsIndex(elastic.AuthorshipMapping())

		wg.Add(7)
		// 1.people
		go func() {
			defer wg.Done()
			people := psql.RetrieveType("Person")
			elastic.AddPeople(people)
		}()
		// 2. affilations
		go func() {
			defer wg.Done()
			affiliations := psql.RetrieveType("Affiliation")
			elastic.AddAffiliationsToPeople(affiliations)
		}()
		// 3. educations
		go func() {
			defer wg.Done()
			educations := psql.RetrieveType("Education")
			elastic.AddEducationsToPeople(educations)
		}()
		// 4. grants
		go func() {
			defer wg.Done()
			grants := psql.RetrieveType("Grant")
			elastic.AddGrants(grants)
		}()
		// 5. funding-roles
		go func() {
			defer wg.Done()
			roles := psql.RetrieveType("FundingRole")
			elastic.AddFundingRoles(roles)
		}()
		// 6. publications
		go func() {
			defer wg.Done()
			publications := psql.RetrieveType("Publication")
			elastic.AddPublications(publications)
		}()
		// 7. authorships
		go func() {
			defer wg.Done()
			authorships := psql.RetrieveType("Authorship")
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
			preview(*typeName)
		} else {
			persistResources(*typeName)
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
