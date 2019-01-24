package graphql_endpoint_test

import (
	"fmt"
	"github.com/OIT-ads-web/graphql_endpoint"
	"testing"
)

func TestPagingPerPage(t *testing.T) {
	var start = 1
	var perPage = 10
	var total = 44

	pageInfo1 := graphql_endpoint.FigurePaging(start, perPage, total)
	if pageInfo1.PerPage != 10 {
		t.Error(fmt.Printf("should be page 10 per page if perPage = %d\n", perPage))
	}
}

func TestPagingTotal(t *testing.T) {
	var start = 1
	var perPage = 10
	var total = 44

	pageInfo1 := graphql_endpoint.FigurePaging(start, perPage, total)
	if pageInfo1.TotalPages != 5 {
		t.Error(fmt.Printf("should be page 5 pages with size=%d and count=%d\n", perPage, total))
	}
}

func TestPagingStart(t *testing.T) {
	var start = 1
	var perPage = 10
	var total = 44

	pageInfo1 := graphql_endpoint.FigurePaging(start, perPage, total)
	if pageInfo1.CurrentPage != 1 {
		t.Error(fmt.Printf("should be page 1 if start = %d and perPage = %d", start, perPage))
	}
}

func TestPagingSecond(t *testing.T) {
	var start = 14
	var perPage = 10
	var total = 44

	pageInfo1 := graphql_endpoint.FigurePaging(start, perPage, total)
	if pageInfo1.CurrentPage != 2 {
		t.Error(fmt.Printf("should be page 2 if start = %d and perPage = %d", start, perPage))
	}
}

func TestPagingOnePast(t *testing.T) {
	var start = 31
	var perPage = 10
	var total = 44

	pageInfo1 := graphql_endpoint.FigurePaging(start, perPage, total)
	if pageInfo1.CurrentPage != 4 {
		t.Error(fmt.Printf("should be page 4 if start = %d and perPage = %d", start, perPage))
	}
}
