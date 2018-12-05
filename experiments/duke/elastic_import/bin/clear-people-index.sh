#!/bin/sh
#curl -XDELETE localhost:9200/people

# also
cd ../
cmd/elastic_import/elastic_import -remove=true -type=people

