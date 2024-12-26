package api

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/skyterra/elastic-embedding-searcher/elastic"
	pb "github.com/skyterra/elastic-embedding-searcher/pb/searcher"
	"github.com/skyterra/elastic-embedding-searcher/runner"
	"io"
)

const IndexReset = "reset"

// IndexCsvFile indexes CSV file data into Elasticsearch.
func (s *SearcherServer) IndexCsvFile(ctx context.Context, req *pb.IndexCsvRequest) (*pb.IndexCsvResponse, error) {
	// decode csv data with base64 and parse total row count.
	data, rowCount, err := DecodeCsvData(req.FileData)
	if err != nil {
		return nil, err
	}

	if rowCount <= 0 {
		return nil, errors.New("csv hs NO data")
	}

	// read csv headers.
	reader, headers, err := ReadCsvHeader(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read csv. err:%s", err.Error())
	}

	embeddingColumn, exist := headers[req.EmbeddingColumn]
	if !exist {
		return nil, fmt.Errorf("no exist embedding column. EmbeddingColumn:%s", req.EmbeddingColumn)
	}

	// create documents, annotations and all keywords(each row).
	documents := make([]*elastic.Document, 0, rowCount)
	annotations := make([]string, 0, rowCount)

	id := 0
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read csv data. err:%s", err.Error())
		}

		id++
		doc := &elastic.Document{
			BaseDocument: elastic.BaseDocument{ID: fmt.Sprintf("%d", id)},
			Metadata:     make(map[string]interface{}),
		}

		for colName, colPos := range headers {
			doc.Metadata[colName] = row[colPos]
		}

		documents = append(documents, doc)
		annotations = append(annotations, row[embeddingColumn])
	}

	// generate embedding vector and save to ES.
	dimsPart1, dimsPart2 := 0, 0
	if len(annotations) > 0 {
		embeddings, err := runner.GenEmbeddingList(ctx, annotations)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding. err:%s", err.Error())
		}

		dimsPart1 = len(embeddings[0]) / 2
		dimsPart2 = len(embeddings[0]) - dimsPart1

		for i := 0; i < len(documents); i++ {
			documents[i].EmbeddingPart1 = embeddings[i][:dimsPart1]
			documents[i].EmbeddingPart2 = embeddings[i][dimsPart1:]
		}
	}

	// reset index if you need.
	if req.ResetIndex == IndexReset && len(annotations) > 0 {
		indexMapping := fmt.Sprintf(elastic.IndexEmbeddingMapping, dimsPart1, dimsPart2)
		if err = elastic.ResetIndex(ctx, req.IndexName, indexMapping); err != nil {
			return nil, err
		}
	}

	// sync to elastic-search.
	docs := make([]elastic.IDocument, 0, len(documents))
	for _, doc := range documents {
		docs = append(docs, doc)
	}

	err = elastic.IndexDocuments(ctx, req.IndexName, docs)
	if err != nil {
		return nil, err
	}

	return &pb.IndexCsvResponse{Message: "OK"}, nil
}

// DecodeCsvData decodes base64 encoded CSV data and returns the decoded data and row count.
func DecodeCsvData(data []byte) ([]byte, int, error) {
	// remove BOM.
	bom := []byte{0xEF, 0xBB, 0xBF}
	if bytes.HasPrefix(data, bom) {
		data = data[len(bom):]
	}

	// remove 'NUL' if existed.
	data = bytes.TrimRight(data, "\x00")

	// check to see if data is empty.
	if len(data) == 0 {
		return nil, 0, nil
	}

	count := bytes.Count(data, []byte{'\n'}) + 1
	return data, count, nil
}

// ReadCsvHeader creates a CSV reader and returns the headers as a map with column positions.
func ReadCsvHeader(data []byte) (*csv.Reader, map[string]int, error) {
	reader := csv.NewReader(bytes.NewBuffer(data))

	// the first line is the header row.
	headers, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}

	// store header name and column position in a map.
	m := make(map[string]int)
	for col, name := range headers {
		m[name] = col
	}

	return reader, m, nil
}
