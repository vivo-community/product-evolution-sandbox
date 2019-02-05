# Search Site

Just a very rudimentary search site of existing elastic index

## run:
npm run start --scripts-prepend-node-path 

## ideas:

### elastic + `ALL_TEXT` ???

Add a field:

> `curl -XPUT 'http://localhost:9200/twitter/_mapping/_doc' -d '`

```json
{
  "properties": {
    "ALL_TEXT": {
      "type": "keyword"
    }
  }
}
```

or `copy_to`:
https://www.elastic.co/guide/en/elasticsearch/reference/current/copy-to.html

### change feed ?
./plugin install https://github.com/jurgc11/es-change-feed-plugin/releases/download/{version}/es-changes-feed-plugin.zip

or

https://dzone.com/articles/elasticsearch5-how-to-build-a-plugin-and-add-a-lis

or logstash:

https://stackoverflow.com/questions/46592747/elasticsearch-to-kafka-event-on-each-change-using-logstash
https://codeforgeek.com/2017/10/elasticsearch-change-feed/



