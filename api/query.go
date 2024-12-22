package api

import (
	"context"
	"errors"
	"github.com/skyterra/elastic-embedding-searcher/elastic"
	"github.com/skyterra/elastic-embedding-searcher/modelx_runner"
	pb "github.com/skyterra/elastic-embedding-searcher/pb/searcher"
	"strings"
)

// Query query document base semantic search .
func (s *SearcherServer) Query(ctx context.Context, req *pb.QueryRequest) (*pb.QueryResponse, error) {
	text := strings.TrimSpace(req.QueryText)
	if text == "" {
		return nil, errors.New("query is empty")
	}

	embedding, err := modelx_runner.GenEmbedding(ctx, text)
	if err != nil {
		return nil, err
	}

	var script string
	if len(embedding) > elastic.MaxNumberOfDimensions {
		part1 := modelx_runner.EmbeddingToString(embedding[:elastic.MaxNumberOfDimensions])
		part2 := modelx_runner.EmbeddingToString(embedding[elastic.MaxNumberOfDimensions:])
		script = elastic.BuildQueryByEmbedding(part1, part2, req.Size)
	} else {
		data := modelx_runner.EmbeddingToString(embedding)
		script = elastic.BuildQueryByEmbedding(data, "", req.Size)
	}

	records, err := elastic.Query(ctx, req.IndexName, script)
	if err != nil {
		return nil, err
	}

	resp := &pb.QueryResponse{
		Message: "succeed",
		Records: make([]*pb.Record, 0, len(records)),
	}

	for _, record := range records {
		r := &pb.Record{
			Id:    record.ID,
			Score: record.Score,
		}

		for k, v := range record.Source.Metadata {
			r.Metadata[k] = v.(string)
		}

		resp.Records = append(resp.Records, r)
	}

	return resp, nil
}
