package elastic

const IndexEmbeddingMapping = `
{
    "settings": {
        "number_of_shards": 2,
        "number_of_replicas": 1
    },
    "mappings": {
        "dynamic": "true",
        "_source": {
            "enabled": "true"
        },
        "properties": {
            "embedding_part1": {
                "type": "dense_vector",
                "dims": %d
            },
            "embedding_part2": {
                "type": "dense_vector",
                "dims": %d
            },
			"metadata": {
				"type": "object",
				"dynamic": "true",
        		"properties": {
          			"keywords": {
            			"type": "keyword" 
          			}
        		}
			}
        }
    }
}`

const BulkIndexTemplate = `{"index": {"_index":"%s", "_id":"%s"}}`
const BulkDelIndexTemplate = `{"delete": {"_index":"%s", "_id":"%s"}}`

const QueryNormalTemplate = `
{
	"query": %s,
	"size": %d,	
	"collapse": {"field": "title.keyword"}
}
`

const QueryEmbeddingTemplate = `
{
    "script_score": {
        "query": {
            "match_all": {

            }
        },
        "script": {
            "source": "cosineSimilarity(params.query_vector_part1, 'embedding_part1') + 1.0 + cosineSimilarity(params.query_vector_part2, 'embedding_part2') + 1.0",
            "params": {
                "query_vector_part1": [%s],
				"query_vector_part2": [%s]
            }
        }
    }
}
`
