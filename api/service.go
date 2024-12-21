package api

import (
	"context"
	pb "github.com/skyterra/elastic-embedding-searcher/pb/searcher"
)

type SearcherServer struct {
	pb.UnimplementedElasticEmbeddingSearcherApiServer
}

func (s *SearcherServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PongResponse, error) {
	return &pb.PongResponse{Status: "Pong"}, nil
}
