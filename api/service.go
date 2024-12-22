package api

import (
	"context"
	pb "github.com/skyterra/elastic-embedding-searcher/pb/searcher"
)

type SearcherServer struct {
	pb.UnimplementedElasticEmbeddingSearcherApiServer
}

func (s *SearcherServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong"}, nil
}
