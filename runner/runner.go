package runner

import (
	"context"
	"errors"
	pb "github.com/skyterra/elastic-embedding-searcher/pb"
	"os"
	"strconv"
	"strings"
)

const (
	uds = "unix:/tmp/grpc_unix_socket_modelx"
)

var ins *ModelXManager

// GetClient returns the ModelX client or panics if not initialized.
func GetClient() pb.ModelxClient {
	if ins == nil {
		panic("not invoke StartModelX")
	}

	ins.clilock.RLock()
	client := ins.client
	ins.clilock.RUnlock()

	return client
}

// StartModelX initializes and starts the ModelX service.
func StartModelX(workers int, cmdPath, modelPath string) error {
	if ins != nil {
		return errors.New("ModelX service has been already started")
	}

	if _, err := os.Stat(cmdPath); err != nil {
		return err
	}

	ins = &ModelXManager{
		stop:      make(chan struct{}),
		workers:   workers,
		cmdPath:   cmdPath,
		modelPath: modelPath,
	}

	if err := ins.fork(); err != nil {
		return err
	}

	if err := ins.dial(); err != nil {
		return err
	}

	return ins.monitor()
}

// StopModelX stops the ModelX service and cleans up.
func StopModelX() error {
	if ins == nil {
		return errors.New("ModelX service not start")
	}

	if err := ins.kill(); err != nil {
		return err
	}

	close(ins.stop)

	ins = nil
	return nil
}

// ParseEmbeddingVector converts a comma-separated string of numbers into a slice of float64.
func ParseEmbeddingVector(content string) ([]float64, error) {
	if len(content) == 0 {
		return nil, errors.New("embedding vector is nil")
	}

	segments := strings.Split(content, ",")
	vector := make([]float64, 0, len(segments))

	for _, s := range segments {
		num, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}

		vector = append(vector, num)
	}

	return vector, nil
}

// ParseEmbeddingVectors converts a slice of comma-separated strings into a slice of float64 slices.
func ParseEmbeddingVectors(contents []string) ([][]float64, error) {
	if len(contents) == 0 {
		return nil, errors.New("embedding vector is nil")
	}

	vectors := make([][]float64, 0, len(contents))
	for _, content := range contents {
		vector, err := ParseEmbeddingVector(content)
		if err != nil {
			return nil, err
		}

		vectors = append(vectors, vector)
	}

	return vectors, nil
}

// GenEmbeddingList generate embedding vectors by passing model.
func GenEmbeddingList(ctx context.Context, annotations []string) ([][]float64, error) {
	reply, err := GetClient().GenEmbeddingList(ctx, &pb.EmbeddingListRequest{TextList: annotations})
	if err != nil {
		return nil, err
	}

	embeddingList, err := ParseEmbeddingVectors(reply.EmbeddingList)
	if err != nil {
		return nil, err
	}

	return embeddingList, nil
}

// GenEmbedding generate embedding vector.
func GenEmbedding(ctx context.Context, text string) ([]float64, error) {
	reply, err := GetClient().GenEmbedding(ctx, &pb.EmbeddingRequest{Text: text})
	if err != nil {
		return nil, err
	}

	embedding, err := ParseEmbeddingVector(reply.Embedding)
	if err != nil {
		return nil, err
	}

	return embedding, nil
}

// CalcSimilarityScore calc similarity score about sourceText and targetTexts.
func CalcSimilarityScore(ctx context.Context, modelName string, sourceText string, targetTexts []string) ([]float64, error) {
	reply, err := GetClient().CalcSimilarityScore(ctx, &pb.SimilarityRequest{
		SourceText:  sourceText,
		TargetTexts: targetTexts,
	})

	if err != nil {
		return nil, err
	}

	return reply.Scores, nil
}

// EmbeddingToString converts a slice of float64 embeddings to a comma-separated string.
func EmbeddingToString(input []float64) string {
	if len(input) == 0 {
		return ""
	}

	result := strconv.FormatFloat(input[0], 'f', -1, 64)
	for _, v := range input[1:] {
		result += "," + strconv.FormatFloat(v, 'f', -1, 64)
	}

	return result
}
