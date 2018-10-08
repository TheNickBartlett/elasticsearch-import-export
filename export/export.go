package main

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var host = flag.String("host", "http://localhost:9200", "Fully qualified Elasticsearch host and port")
var index = flag.String("index", "documents", "The index (or alias) name to run the query against")
var output = flag.String("output", "./output.json.gz", "Output file to store results to")
var query = flag.String("query", "", "Elasticsearch search query JSON")

var outputFile *os.File
var outputFileWriter *gzip.Writer

type scrollResult struct {
	ScrollID string `json:"_scroll_id"`
	Hits     scrollHits
}

type scrollHits struct {
	Total int
	Hits  []scrollHit
}

type scrollHit struct {
	Source json.RawMessage `json:"_source"`
}

type scrollQuery struct {
	Size  int             `json:"size"`
	Query json.RawMessage `json:"query,omitempty"`
}

type nextScrollRequest struct {
	Scroll   string `json:"scroll"`
	ScrollID string `json:"scroll_id"`
}

func main() {
	flag.Usage = printUsage
	flag.Parse()

	checkConnectivity()

	createOutputFileWriter()

	defer outputFile.Close()
	defer outputFileWriter.Close()

	handleScrollQuery()
}

func printUsage() {
	fmt.Println("Overview")
	fmt.Println("Queries an Elasticsearch instance using a scroll query and outputs the source of each result on a new line in a gzipped file")
	fmt.Println()
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])

	flag.PrintDefaults()
}

func checkConnectivity() {
	resp, err := http.Get(*host)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
}

func createOutputFileWriter() {
	var err error
	outputFile, err = os.Create(*output)
	if err != nil {
		panic(err)
	}
	outputFileWriter, _ = gzip.NewWriterLevel(outputFile, gzip.BestSpeed)
}

func handleScrollQuery() {
	scrollResult := doFirstScrollRequest()
	numResults := len(scrollResult.Hits.Hits)

	if numResults > 0 {
		outputHits(scrollResult.Hits.Hits)
	}

	for numResults > 0 {
		scrollResult := doSubsequentScrollRequest(scrollResult.ScrollID)
		numResults = len(scrollResult.Hits.Hits)
		outputHits(scrollResult.Hits.Hits)
	}

}

func doFirstScrollRequest() scrollResult {

	qry := &json.RawMessage{}
	qry.UnmarshalJSON([]byte(*query))

	req := &scrollQuery{
		Size:  1000,
		Query: *qry,
	}
	reqString, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("%s/%s/_search?scroll=1m", *host, *index)

	return makeScrollRequest(string(reqString), url)
}

func doSubsequentScrollRequest(scrollID string) scrollResult {
	req := &nextScrollRequest{
		Scroll:   "1m",
		ScrollID: scrollID,
	}

	reqString, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	url := fmt.Sprintf("%s/_search/scroll", *host)

	return makeScrollRequest(string(reqString), url)
}

func makeScrollRequest(reqBody string, url string) scrollResult {

	reqBodyReader := strings.NewReader(reqBody)
	resp, err := http.Post(url, "application/json", reqBodyReader)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic(errors.New(string(respBody)))
	}

	var res scrollResult
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		panic(err)
	}

	return res
}

func outputHits(hits []scrollHit) {
	for _, hit := range hits {
		outputFileWriter.Write(hit.Source)
		outputFileWriter.Write([]byte("\n"))
	}
}
