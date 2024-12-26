package api

import (
	"context"
	"errors"
	"github.com/skyterra/elastic-embedding-searcher/elastic"
	pb "github.com/skyterra/elastic-embedding-searcher/pb/searcher"
	"github.com/skyterra/elastic-embedding-searcher/runner"
	"strings"
)

// Query query document base semantic search .
func (s *SearcherServer) Query(ctx context.Context, req *pb.QueryRequest) (*pb.QueryResponse, error) {
	text := strings.TrimSpace(req.QueryText)
	if text == "" {
		return nil, errors.New("query is empty")
	}

	embedding, err := runner.GenEmbedding(ctx, text)
	if err != nil {
		return nil, err
	}

	part1 := runner.EmbeddingToString(embedding[:len(embedding)/2])
	part2 := runner.EmbeddingToString(embedding[len(embedding)/2:])
	script := elastic.BuildQueryByEmbedding(part1, part2, req.Size)

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
			Id:       record.ID,
			Score:    record.Score,
			Metadata: make(map[string]string),
		}

		for k, v := range record.Source.Metadata {
			r.Metadata[k] = v.(string)
		}

		resp.Records = append(resp.Records, r)
	}

	return resp, nil
}
