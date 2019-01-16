# graphql server in golang

## config 

* see config.toml.example

## server 

* endpoint on :9001 (not configurable at this point)
* see localhost:9001/graphql


## cors - (test)
curl -H "Origin: http://example.com" \
                       -H "Access-Control-Request-Method: POST" \
                       -H "Access-Control-Request-Headers: X-Requested-With" \
                       -X OPTIONS --verbose http://localhost:9001/graphql
