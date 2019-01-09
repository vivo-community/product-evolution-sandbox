#!/bin/sh

cd cmd/elastic_query
go build
cd ../../
cd cmd/graphql_server
go build


