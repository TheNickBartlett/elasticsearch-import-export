# Overview (export.go)
Provides a program that queries Elasticsearch for documents based on the given query, and outputs the results to a new line separated json gzip file

Run `export -h` for all command line argument options.

## Flags
### `-host`
The fully qualified Elasticsearch server url. E.g. `https://myAWSElasticsearchDomain`, or `http://localhost:9200`

### `-index`
The name of the index, or alias, that the query should run against

### `-output`
The relative or absolute path to the target output file. E.g. `./output.json.gz`, or `/tmp/export/documents.json.gz`

### `-query`
The Elasticsearch JSON query to use to limit documents, minus the outer `query` object. Note, quotes must be escaped.
E.g. `{\"terms\":{\"type\":[\"a\",\"b\",\"c\"]}}`

## Example Usage
`./export -index="myindex" -query="{\"terms\":{\"type\":[\"a\",\"b\",\"c\"]}}" -host="https://my-aws-es-cluster.region.es.amazonaws.com" -output="./documents.json.gz"` 

# Overview (import.go)
Provides a program that uses a new line separated json gzip file produced by `export` as a data source to import into an Elasticsearch index.

Run `import -h` for all command line argument options.

## Flags
### `-host`
The fully qualified Elasticsearch server url. E.g. `https://myAWSElasticsearchDomain`, or `http://localhost:9200`

### `-file`
The relative or absolute path to the source file. E.g. `./output.json.gz`, or `/tmp/export/documents.json.gz`

### `-index`
The name of the index, or alias, that the data source should import to

### `-type`
The Elasticsearch document type to use when importing

## Example Usage
`./import -index="documents" -type="doc" -file="./documents.json.gz" -host="https://my-aws-es-cluster.region.es.amazonaws.com"`