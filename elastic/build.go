package elastic

import "fmt"

// BuildQueryByEmbedding builds a search query using large embeddings.
func BuildQueryByEmbedding(embeddingPart1, embeddingPart2 string, size int32) string {
	scriptScore := fmt.Sprintf(QueryEmbeddingTemplate, embeddingPart1, embeddingPart2)
	return fmt.Sprintf(QueryNormalTemplate, scriptScore, size)
}
