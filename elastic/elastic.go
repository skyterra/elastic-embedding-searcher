package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"io"
	"strings"
	"time"
)

var Client *elasticsearch.Client

func Init(address string, username, password string) error {
	// create elastic client.
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: strings.Split(address, ","),
		Username:  username,
		Password:  password,
	})

	if err != nil {
		return err
	}

	Client = client
	return nil
}

func Cleanup() error {
	return nil
}

// ResetIndex deletes and recreates an index with the given name and mapping.
func ResetIndex(ctx context.Context, indexName string, indexMapping string) error {
	_, err := Client.Indices.Delete([]string{indexName}, Client.Indices.Delete.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to delete index. index:%s err:%s", indexName, err.Error())
	}

	response, err := Client.Indices.Create(
		indexName,
		Client.Indices.Create.WithBody(strings.NewReader(indexMapping)),
		Client.Indices.Create.WithContext(ctx),
	)

	if err != nil {
		return fmt.Errorf("failed to create index. index:%s err:%s", indexName, err.Error())
	}

	if response.IsError() {
		return fmt.Errorf("failed to create index. index:%s err:%s", indexName, response.String())
	}

	return nil
}

// IndexDocuments indexes multiple documents in the specified Elasticsearch index.
func IndexDocuments(ctx context.Context, indexName string, documents []IDocument) error {
	var buf bytes.Buffer

	for i := 0; i < len(documents); i++ {
		data, err := json.Marshal(documents[i])
		if err != nil {
			continue
		}

		switch documents[i].GetOperation() {
		case BulkDelete:
			buf.Write([]byte(fmt.Sprintf(BulkDelIndexTemplate, indexName, documents[i].GetID())))
			buf.WriteByte('\n')

		default:
			buf.Write([]byte(fmt.Sprintf(BulkIndexTemplate, indexName, documents[i].GetID())))
			buf.WriteByte('\n')

			buf.Write(data)
			buf.WriteByte('\n')
		}
	}

	response, err := Client.Bulk(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.IsError() {
		return errors.New(response.String())
	}

	return nil
}

// Query queries general Document base on elasticsearch.
func Query(ctx context.Context, indexName, script string) ([]QueryRecord, error) {
	return query[Document](ctx, indexName, script)
}

// query performs a search on an Elasticsearch index using a provided script.
func query[DocType TDocument](ctx context.Context, indexName, script string) ([]TQueryRecord[DocType], error) {
	response, err := Client.Search(
		Client.Search.WithContext(ctx),
		Client.Search.WithIndex(indexName),
		Client.Search.WithBody(strings.NewReader(script)),
		Client.Search.WithTrackTotalHits(true),
		Client.Search.WithTimeout(time.Second),
	)

	if err != nil {
		return nil, err
	}

	// I always remember to close the body.
	defer response.Body.Close()

	if response.IsError() {
		data, _ := io.ReadAll(response.Body)
		return nil, errors.New(string(data))
	}

	ss := &TSearchResponse[DocType]{}
	err = json.NewDecoder(response.Body).Decode(ss)
	if err != nil {
		return nil, err
	}

	return ss.Hits.Hits, nil
}
