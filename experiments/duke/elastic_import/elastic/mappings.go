package elastic

import (
	"fmt"
	wi "github.com/OIT-ads-web/widgets_import"
)


func PersonMapping() string {
	fmt.Println("trying to show PersonMapping")
	return wi.RenderTemplate("person.tmpl")
}

func AffiliationMapping() string {
	return wi.RenderTemplate("affiliation.tmpl")
}

func FundingRoleMapping() string {
	return wi.RenderTemplate("funding-role.tmpl")
}

func PublicationMapping() string {
	return wi.RenderTemplate("publication.tmpl")
}

func AuthorshipMapping() string {
	return wi.RenderTemplate("authorship.tmpl")
}

func GrantMapping() string {
	return wi.RenderTemplate("grant.tmpl")
}

