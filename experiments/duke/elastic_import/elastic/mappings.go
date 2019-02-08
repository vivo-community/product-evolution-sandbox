package elastic

import (
	wi "github.com/OIT-ads-web/widgets_import"
)


func PersonMapping() string {
	return wi.RenderTemplate("person.tmpl")
}

func AffiliationMapping() string {
	return wi.RenderTemplate("affiliation.tmpl")
}

func FundingRoleMapping() string {
	return wi.RenderTemplate("funding-role")
}

func PublicationMapping() string {
	return wi.RenderTemplate("publication")
}

func AuthorshipMapping() string {
	return wi.RenderTemplate("authorship")
}

func GrantMapping() string {
	return wi.RenderTemplate("grant")
}

