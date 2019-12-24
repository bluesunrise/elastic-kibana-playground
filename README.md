# elastic-kibana-playground
Playground for playing around with Elastic API with friends

(1) To get this example to work, first open up the docker-compose ports in docker-compose.yml:

Elastic Search, uncomment:
``` 
ports:
      - "9200:9200"
```
Kibana, uncomment:
```
	ports:
      - "5601:5601"
```

(2) Create a `Saved Query` in Kibana - do not create a `Saved Search`, see:

https://www.elastic.co/guide/en/kibana/master/save-load-delete-query.html

(You will first need to setup an index in Kibana)

Here is how I generated output to search-output.json for a Kibana Query (see search-output.json):

 curl http://localhost:5601/kibana/api/saved_objects/_find?type=query

(3) Start writing the code ...

https://www.elastic.co/guide/en/kibana/master/using-api.html

###PROTOTYPE GOAL:
	* retrieve all saved queries (done)
	* use the saved query to execute an Elastic Search *KQL QUERY* as a metric check
	
For example, finding log entries of some search criteria happening over a time period, see search-output.txt
