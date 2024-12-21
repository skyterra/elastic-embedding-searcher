package api

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/skyterra/elastic-embedding-searcher/pb/searcher"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/signal"
	"syscall"
)

func Start(port int16) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := grpc.NewServer()
	pb.RegisterElasticEmbeddingSearcherApiServer(server, &SearcherServer{})

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	go func() {
		log.Printf("start at:%d", port)
		if err := server.Serve(listen); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			panic(err)
		}
	}()

	<-ctx.Done()
	stop()

	server.GracefulStop()
	log.Println("shutting down gracefully.")
	return nil
}

func Cleanup() error {

	return nil
}

func StartService(port int) error {

	return nil
}
