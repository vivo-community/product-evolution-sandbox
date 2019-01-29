package graphql_endpoint_test

import (
	"fmt"
	"github.com/OIT-ads-web/graphql_endpoint"
	"testing"
)

func TestPagingPerPage(t *testing.T) {
	var size = 10
	var from = 1
	var total = 44

	pageInfo1 := graphql_endpoint.FigurePaging(size, from, total)
	if pageInfo1.PerPage != 10 {
		t.Error(fmt.Printf("should be page 10 per page if size = %d\n", size))
	}
}

func TestPagingTotal(t *testing.T) {
	var size = 10
	var from = 1
	var total = 44

	pageInfo1 := graphql_endpoint.FigurePaging(size, from, total)
	if pageInfo1.TotalPages != 5 {
		t.Error(fmt.Printf("should be page 5 pages with size=%d and count=%d\n", size, total))
	}
}

func TestPagingStart(t *testing.T) {
	var size = 10
	var from = 1
	var total = 44

	pageInfo1 := graphql_endpoint.FigurePaging(size, from, total)
	if pageInfo1.CurrentPage != 1 {
		t.Error(fmt.Printf("should be page 1 if start = %d and perPage = %d",from, size))
	}
}

func TestPagingSecond(t *testing.T) {
	var size = 10
	var from = 14
	var total = 44

	pageInfo1 := graphql_endpoint.FigurePaging(size, from, total)
	if pageInfo1.CurrentPage != 2 {
		t.Error(fmt.Printf("should be page 2 if start = %d and perPage = %d", from, size))
	}
}

func TestPagingOnePast(t *testing.T) {
	var size = 10
	var from = 31
	var total = 44

	pageInfo1 := graphql_endpoint.FigurePaging(size, from, total)
	if pageInfo1.CurrentPage != 4 {
		t.Error(fmt.Printf("should be page 4 if from = %d and size = %d", from, size))
	}
}

func TestPagingSmall(t *testing.T) {
	var size = 10
	var from = 1
	var total = 3

	pageInfo1 := graphql_endpoint.FigurePaging(size, from, total)
	if pageInfo1.CurrentPage != 1 {
		t.Error(fmt.Printf("should be page 1 if from = %d and size = %d", from, size))
	}
	if pageInfo1.TotalPages != 1 {
		t.Error(fmt.Printf("should be 1 pages(s) if from = %d and size = %d", from, size))
	}
}
