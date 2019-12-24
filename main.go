package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	storedQueries, _ := retrieveStoredQueries();
	storedQuery := storedQueries[0]
	log.Print(storedQuery)

	// TODO: use storedQuery variable in place of *INSERT_QUERY_HERE*

	// TODO: these should come from environment variables
	cfg := elasticsearch.Config{

		Addresses: []string{ // TODO: multiple load balanced elastic search
			"http://localhost:9200",
			// "http://localhost:9201",
		},
	}
	esClient, _ := elasticsearch.NewClient(cfg)
	log.Println(elasticsearch.Version)
	log.Println(esClient.Info())

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"container.name": "dockergw8_nagios2collage_eventlog_1", // *INSERT_QUERY_HERE*
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	var (
		r  map[string]interface{}
	//	wg sync.WaitGroup
	)
	// Perform the search request.
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex("logstash-*"),
		esClient.Search.WithBody(&buf),
		esClient.Search.WithTrackTotalHits(true),
		esClient.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	log.Println(strings.Repeat("=", 37))
	defer res.Body.Close()
}

func retrieveStoredQueries() ([]string, int) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := http.Client{Transport: tr}
	var request *http.Request
	var response *http.Response
	var err error

	// types of saved objects -- Valid options include QUERY visualization, dashboard, search, index-pattern, config, and timelion-sheet.
	// kibanaFilter := strings.NewReader("{ \"type\": \"index-pattern\", \"excludeExportDetails\": true }")

	request, err = http.NewRequest(http.MethodGet, "http://localhost:5601/kibana/api/saved_objects/_find?type=query", nil);
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("kbn-xsrf", "true")

	response, err = client.Do(request)
	if response.StatusCode == 400 {
		log.Fatalf("Not Found!")
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	result := string(responseBody)
	response.Body.Close()
	// TODO: parse out the query section of the JSON and return it
	return []string{result}, 0
}