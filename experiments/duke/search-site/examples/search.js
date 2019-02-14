import { client } from './config'

var inspect = function (resp) {
    var hits = resp.hits.hits;
    var len = hits.length;
    for (var i = 0; i < len; i++) {
      var result = hits[i]._source
      console.log(result);
      console.log(result.name);
      console.log(result.primaryTitle);
      console.log(result.type);
    }
};

var inspectAgg = function (resp) {
    console.log("*****************")
    console.log(resp)
    console.log(resp.aggregations.buckets)
    var hits = resp.hits.hits;
    var len = hits.length;
    for (var i = 0; i < len; i++) {
      //var result = hits[i]._source
      //console.log(result);
      //console.log(result.name);
      //console.log(result.primaryTitle);
      //console.log(result.type);
    }
    console.log("******************")
};


var error = function(err) {
  console.log("ERROR")
  console.trace(err.message)
}

export { inspect, error }

/*
client.search({
  index: 'people',
  //type: 'person',
  body: {
    query: {
      match: {
        'name.lastName': 'Swaminathan'
        //'keywordList.label': 'Cardiomyopathies'
        //'_all': 'Swaminathan'
      }
    }
  }
}).then(function (resp) {
    var hits = resp.hits.hits;
    var len = hits.length;
    for (var i = 0; i < len; i++) {
      var result = hits[i]._source
      console.log(result);
      console.log(result.name);
      console.log(result.primaryTitle);
      console.log(result.type);
    }
}, function (err) {
    console.trace(err.message);
});


*/
