# Overview (export.go)
Provides a script that queries Elasticsearch for documents based on the given query, and outputs the results to a new line separated gzip json file

Run `go run export.go -h` for all command line argument options

## Flags
### `-host`
The fully qualified Elasticsearch server url. E.g. `https://myAWSElasticsearchDomain`, or `http://localhost:9200`

### `-index`
The name of the index, or alias, that the query should run against

### `-output`
The fully qualified path to the target output file. E.g. `./output.json.gz`, or `/tmp/export/documents.json.gz`

### `-query`
The Elasticsearch JSON query to use to limit documents. 
E.g. `{"terms":{"type":["a","b","c"]}}`

## Example Usage
`go run export.go -query="{\"terms\":{\"type\":[\"a\",\"b\",\"c\"]}}"`

# Overview (import.go)
Soon(tm)