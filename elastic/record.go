package elastic

// TDocument define supporting document types.
type TDocument interface {
	Document
}

// TQueryRecord represents a single search result record.
type TQueryRecord[DocType TDocument] struct {
	Source DocType `json:"_source"`
	Score  float64 `json:"_score"`
	ID     string  `json:"_id"`
}

// THits represents a collection of search result records.
type THits[DocType TDocument] struct {
	Hits []TQueryRecord[DocType] `json:"hits"`
}

// TSearchResponse represents the response from a search query.
type TSearchResponse[DocType TDocument] struct {
	Hits THits[DocType] `json:"hits"`
}

// TBulkSearchResponse represents a response for bulk search operations.
type TBulkSearchResponse[DocType TDocument] struct {
	Responses []TSearchResponse[DocType] `json:"responses"`
}

// QueryRecord represents general Document record.
type QueryRecord = TQueryRecord[Document]

// BulkSearchResponse represents a response structure(general Document) for bulk search operations.
type BulkSearchResponse = TBulkSearchResponse[Document]
