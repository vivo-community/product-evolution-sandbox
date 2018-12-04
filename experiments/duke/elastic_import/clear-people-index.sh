#!/bin/sh
curl -XDELETE localhost:9200/people

# also
# cmd/elastic_import/elastic_import -remove -type people

