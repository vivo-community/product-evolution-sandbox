#!/bin/sh
#curl -XDELETE localhost:9200/affiliations

# also
cd ../
cmd/elastic_import/elastic_import -remove=true -type=affiliations

