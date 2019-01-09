#!/bin/sh

cd cmd/widgets_import
go build
cd ../../
cd cmd/elastic_import
go build
cd ../../
cd cmd/staging_import
go build



