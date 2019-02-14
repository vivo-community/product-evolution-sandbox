//var host = 'https://elasticsearch-ads-graphql-elastic.cloud.duke.edu';
var host = 'http://localhost:9200';

var elasticsearch = require('elasticsearch');
var client = new elasticsearch.Client({
  host: host,
  log: 'trace'
});

export { client };
