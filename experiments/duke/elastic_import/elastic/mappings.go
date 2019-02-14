package elastic

import (
	wi "github.com/OIT-ads-web/widgets_import"
)

// make these return error?  
func PersonMapping() (string, error) {
	return wi.RenderTemplate("person.tmpl")
}

func AffiliationMapping() (string, error) {
	return wi.RenderTemplate("affiliation.tmpl")
}

func FundingRoleMapping() (string, error) {
	return wi.RenderTemplate("funding-role.tmpl")
}

func PublicationMapping() (string, error) {
	return wi.RenderTemplate("publication.tmpl")
}

func AuthorshipMapping() (string, error) {
	return wi.RenderTemplate("authorship.tmpl")
}

func GrantMapping() (string, error) {
	return wi.RenderTemplate("grant.tmpl")
}

