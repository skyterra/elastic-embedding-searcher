package elastic

type BulkOperation uint8

const (
	BulkIndex  BulkOperation = 0
	BulkDelete BulkOperation = 1

	MaxNumberOfDimensions = 2048
)

type IDocument interface {
	GetID() string
	GetOperation() BulkOperation
}

type BaseDocument struct {
	ID        string        `json:"ID"`
	Operation BulkOperation `json:"-"`
}

// Document represents a document structure with metadata and large embedding information.
//
// large embedding: the dimension of embedding vector greater than MaxNumberOfDimensions.
//
// For large document queries, the Document structure is typically used to return results.
// From a business perspective, the embedding-related fields are generally ignored in the response.
type Document struct {
	BaseDocument

	EmbeddingPart1 []float64              `json:"embedding_part1"`
	EmbeddingPart2 []float64              `json:"embedding_part2"`
	Metadata       map[string]interface{} `json:"metadata"`
}

func (doc *BaseDocument) GetID() string {
	return doc.ID
}

func (doc *BaseDocument) GetOperation() BulkOperation {
	return doc.Operation
}
