package emb

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
func GenEmbeddingList(ctx context.Context, modelName string, annotations []string) ([][]float64, error) {
	reply, err := GetClient().GenEmbeddingList(ctx, &pb.EmbeddingListRequest{TextList: annotations, ModelName: modelName})
	if err != nil {
		return nil, err
	}

	embeddingList, err := ParseEmbeddingVectors(reply.EmbeddingList)
	if err != nil {
		return nil, err
	}

	return embeddingList, nil
}

// GenEmbedding generate embedding vector by passing model.
func GenEmbedding(ctx context.Context, modelName string, text string) ([]float64, error) {
	reply, err := GetClient().GenEmbedding(ctx, &pb.EmbeddingRequest{Text: text, ModelName: modelName})
	if err != nil {
		return nil, err
	}

	embedding, err := ParseEmbeddingVector(reply.Embedding)
	if err != nil {
		return nil, err
	}

	return embedding, nil
}

// GenStringEmbedding generates an embedding vector for a given text and index name.
func GenStringEmbedding(ctx context.Context, modelName string, text string) (string, error) {
	reply, err := GetClient().GenEmbedding(ctx, &pb.EmbeddingRequest{Text: text, ModelName: modelName})
	if err != nil {
		return "", err
	}

	return reply.Embedding, nil
}

// CalcSimilarityScore calc similarity score about sourceText and targetTexts.
func CalcSimilarityScore(ctx context.Context, modelName string, sourceText string, targetTexts []string) ([]float64, error) {
	reply, err := GetClient().CalcSimilarityScore(ctx, &pb.SimilarityRequest{
		SourceText:  sourceText,
		TargetTexts: targetTexts,
		ModelName:   modelName,
	})

	if err != nil {
		return nil, err
	}

	return reply.Scores, nil
}

func SplitSubSentences(annotation string, maxWordCount int, sep string) []string {
	annotation = strings.Trim(annotation, sep)
	words := strings.Split(annotation, sep)

	// words count match requirement.
	if len(words) <= maxWordCount {
		return []string{annotation}
	}

	var subSentences []string

	// group words by maxWordCount.
	for i := 0; i <= len(words)-maxWordCount; i += maxWordCount {
		subSentences = append(subSentences, strings.Join(words[i:i+maxWordCount], sep))
	}

	return subSentences
}
