Imports vivo_widgets data into an Elastic instance

# Go 

This needs go >= 1.11.1 - because it uses the new 'mod' stuff

I've installed that locally with asdf, then set local GOPATH with direnv
but eventually this will be dockerized

# Setup

Requires a bit (now):

## postgresql
Running instance of postgresql with credentials that are
read from a './config.toml' file (see config.toml.example):
	
## elastic
Running instance of elasticsearch
* localhost:9200 

NOTE: these can both be started via docker-compose

> docker-compose up

or if one or other is already install locally

> docker-compose up <service_to_start>

## Importers (e.g. Go Code)

### Building

> ./build.sh

### Running

> cmd/widgets_import/widgets_import -config <config_file> -dry-run (true|false)
> cmd/elastic_import/elastic_import -config <config_file> 

[NOTE: only run once, unless -remove=true]

If you leave off the -config option it will default to ./config.toml

To delete elastic index (people for instance):
> curl -XDELETE localhost:9200/people

or

> cmd/elastic_import/elastic_import -type=people -remove=true


### Elastic caveats

It is possible that elastic_import may only work inside of the Docker instance. In which case you would do the following:

> cd produce-evolution-sandbox/experiments/duke/elastic_import
> docker-compose run --rm importer sh
> bash
> ./build.sh
> ./cmd/elastic_import/elastic_import -type all

If localhost:9200 does not display the Elasticsearch instance in your browser, try docker:9200 instead.


### Getting started

#### People
> cmd/widgets_import/widgets_import -type=people
> cmd/elastic_import/elastic_import -type=people

#### Affiliations
> cmd/widgets_import/widgets_import -type=positions
> cmd/elastic_import/elastic_import -type=affiliations

#### Educations
> cmd/widgets_import/widgets_import -type=educations
> cmd/elastic_import/elastic_import -type=educations

#### Grants
> cmd/widgets_import/widgets_import -type=grants
> cmd/elastic_import/elastic_import -type=grants

#### Publications 
> cmd/widgets_import/widgets_import -type=publications
> cmd/elastic_import/elastic_import -type=publications


## Exporting data from Elasticsearch
> npm install elasticdump

### Backup index map to a file:
> ./bin/elasticdump \
  --input=http://docker:9200/people \
  --output=people_index_mapping.json \
  --type=mapping

### Backup index data to a file:
> ./bin/elasticdump \
  --input=http://docker:9200/people \
  --output=people_index_data.json \
  --type=data

For additional elasticdump commands and options, see: https://hub.docker.com/r/taskrabbit/elasticsearch-dump/

## Elastic Mappings

If you want to use the graphql_endpoint and react-static parts of this duke/experiment
folder without using this particular, duke-specific data ingest method, you only need
to bring in data however it is easiest for you in the following elastic mappings:


### personMapping

```json
"person":{
	"properties":{
		"id":           { "type": "text" },
		"uri":          { "type": "text" },
		"primaryTitle": { "type": "text" },
		"name":{
			"type":"object",
			"properties": {
				"firstName":  { "type": "text" },
				"lastName":   { "type": "text" },
				"middleName": { "type": "text" }
		    }
		},
		"image": {
			"type": "object",
			"properties": {
				"main":      { "type": "text" },
				"thumbnail": { "type": "text" }
			}
		},
	    "keywordList": {
	      "type": "nested",
	      "properties": {
		      "uri":   { "type": "text" },
		      "label": { "type": "text" }
	      }
	    }
    }
}
```

### affiliations (positions)

```json
"affiliation":{
	"properties":{
		"id":        { "type": "text" },
		"uri":       { "type": "text" },
		"personId":  { "type": "text" },
		"label":     { "type": "text" },
		"startDate": {
			"type": "object",
			"properties": {
				"dateTime":   { "type": "text" },
				"resolution": { "type": "text" }
			}
		},
		"organizationId":    { "type": "text" },
		"organizationLabel": { "type": "text" } 
    }
}
```

### educations

```json
"education":{
	"properties":{
		"id":        { "type": "text" },
		"uri":       { "type": "text" },
		"label":     { "type": "text" },
		"personId":  { "type": "text" },
		"org":     { 
			"type": "object",
			"properties": {
				"id": { "type": "text" },
				"label": { "type": "text" }
			}
		}
	}
}
```

### grants

```json
"grant":{
	"properties":{
		"id":        { "type": "text" },
		"uri":       { "type": "text" },
		"label":     { "type": "text" },
		"startDate": {
			"type": "object",
			"properties": {
				"dateTime":   { "type": "text" },
				"resolution": { "type": "text" }
			}
		},
		"endDate": {
			"type": "object",
			"properties": {
				"dateTime":   { "type": "text" },
				"resolution": { "type": "text" }
			}
		}
	}
}
```

### funding-roles

```json
"funding-role":{
	"properties":{
		"id":        { "type": "text" },
		"uri":       { "type": "text" },
		"grantId":   { "type": "text" },
		"personId":  { "type": "text" },
		"label":     { "type": "text" }
	}
}


### publications

```json
"publication":{
	"properties":{
		"id":         { "type": "text" },
		"uri":        { "type": "text" },
		"label":      { "type": "text" },
		"authorList": { "type": "text" },
		"doi":        { "type": "text" },
        "venue":      { 
			"type": "object",
			"properties": {
				"uri":   { "type": "text" },
				"label": { "type": "text" }
			}
		}
	}
}
```

### authorships

```json
"authorship":{
	"properties":{
		"id":             { "type": "text" },
		"uri":            { "type": "text" },
		"publicationId":  { "type": "text" },
		"personId":       { "type": "text" },
		"label":          { "type": "text" }
	}
}
```


