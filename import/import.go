package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var host = flag.String("host", "http://localhost:9200", "Fully qualified Elasticsearch host and port")
var importFileName = flag.String("file", "", "specific json.gz file name containing documents to import. Can be relative or absolute.")
var index = flag.String("index", "", "the name of the index documents will be created in")
var docType = flag.String("type", "", "the name of the type to give to new documents when indexing them")

type bulkIndexRequestMetadata struct {
	Index   string `json:"_index"`
	DocType string `json:"_type"`
	ID      string `json:"_id"`
}

type bulkIndexRequestItem struct {
	Index bulkIndexRequestMetadata `json:"index"`
}

func main() {
	flag.Usage = printUsage
	flag.Parse()

	if *index == "" {
		log.Fatalln("Index name is required")
	}

	if *docType == "" {
		log.Fatalln("type is required")
	}

	if *importFileName == "" {
		log.Fatalln("importFileName is required")
	}

	file, err := os.Open(*importFileName)
	if err != nil {
		log.Fatalf("Unable to open file %s", *importFileName)
	}
	reader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatalf("Unable to read file %s", *importFileName)
	}

	checkConnectivity()

	scanner := bufio.NewScanner(reader)
	const maxCapacity = 1024 * 1024 * 10
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	var lines int64
	batch := make([]string, 1000)

	for scanner.Scan() {
		lines++
		batch = append(batch, scanner.Text())
		if len(batch) >= 1000 {
			importBatch(batch)
			batch = batch[:0]
		}
	}
	importBatch(batch)

	log.Println(scanner.Err())
	log.Printf("processed %d lines", lines)
}

func printUsage() {
	fmt.Println("Overview")
	fmt.Println("Imports a json.gz file previously exported from export.go.")
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

func importBatch(documents []string) {

	requestBodyBuilder := &strings.Builder{}

	for _, document := range documents {
		if len(document) == 0 {
			continue
		}
		m := map[string]interface{}{}
		err := json.Unmarshal([]byte(document), &m)
		if err != nil {
			fmt.Println(document)
			fmt.Println(err)
		}

		customImportLogic(&m)
		documentID := produceDocumentID(&m, document)

		meta := bulkIndexRequestMetadata{
			Index:   *index,
			DocType: *docType,
			ID:      documentID,
		}
		item := bulkIndexRequestItem{
			Index: meta,
		}

		itemStr, err := json.Marshal(item)
		if err != nil {
			panic(err)
		}

		bodyStr, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}

		requestBodyBuilder.Write(itemStr)
		requestBodyBuilder.WriteString("\n")
		requestBodyBuilder.Write(bodyStr)
		requestBodyBuilder.WriteString("\n")
	}

	// Now write these documents to ES
	url := fmt.Sprintf("%s/_bulk", *host)
	resp, err := http.Post(url, "application/x-ndjson", strings.NewReader(requestBodyBuilder.String()))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s %d \n", resp.Status, resp.StatusCode)
}

func customImportLogic(document *map[string]interface{}) {
	// Insert custom logic to manipulate the resulting document here
}

func produceDocumentID(document *map[string]interface{}, sourceDocument string) string {
	// Insert custom logic to generate a document ID here
	return ""
}
