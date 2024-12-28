package messages

import (
	"context"
	"encoding/json"
	"github.com/skyterra/clog"
	"github.com/skyterra/elastic-embedding-searcher/helper"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/skyterra/elastic-embedding-searcher/elastic"
)

type QuoteMessage struct {
	ID     string `json:"id"`
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

func NewMessage(data string) *KafkaMessage {
	return &KafkaMessage{
		Value: []byte(data),
	}
}

var MessageSet = []IMessage{
	NewMessage(`{"id":"1", "quote":"The only limit to our realization of tomorrow is our doubts of today.", "author":"Franklin D. Roosevelt"}`),
	NewMessage(`{"id":"2", "quote":"In the middle of difficulty lies opportunity.", "author":"Albert Einstein"}`),
	NewMessage(`{"id":"3", "quote":"Success is not final, failure is not fatal: It is the courage to continue that counts.", "author":"Winston Churchill"}`),
	NewMessage(`{"id":"4", "quote":"Be yourself; everyone else is already taken.", "author":"Oscar Wilde"}`),
	NewMessage(`{"id":"5", "quote":"Do what you can, with what you have, where you are.", "author":"Theodore Roosevelt"}`),
	NewMessage(`{"id":"6", "quote":"Not everything that is faced can be changed, but nothing can be changed until it is faced.", "author":"James Baldwin"}`),
	NewMessage(`{"id":"7", "quote":"The best way to predict the future is to create it.", "author":"Peter Drucker"}`),
	NewMessage(`{"id":"8", "quote":"Happiness is not something ready made. It comes from your own actions.", "author":"Dalai Lama"}`),
	NewMessage(`{"id":"9", "quote":"It always seems impossible until it’s done.", "author":"Nelson Mandela"}`),
	NewMessage(`{"id":"10", "quote":"You miss 100% of the shots you don’t take.", "author":"Wayne Gretzky"}`),
}

type MockConsumer struct {
}

func (m *MockConsumer) FetchMessage(ctx context.Context) (IMessage, error) {
	return MessageSet[rand.IntN(len(MessageSet))], nil
}

func (m *MockConsumer) CommitMessages(ctx context.Context, msgs ...IMessage) error {
	return nil
}

func (m *MockConsumer) Close() error {
	return nil
}

func init() {
	clog.SetDefaultOpts(helper.ReadContextTrace, helper.ReadContextModule)
}

func MockParseFunc(message []byte) (elastic.IDocument, error) {
	quote := &QuoteMessage{}
	err := json.Unmarshal(message, quote)
	if err != nil {
		return nil, err
	}

	doc := &elastic.Document{
		BaseDocument: elastic.BaseDocument{ID: quote.ID, Operation: elastic.BulkIndex},
		Metadata: map[string]interface{}{
			"author": quote.Author,
			"quote":  quote.Quote,
		},
	}

	return doc, nil
}

func TestNewMessageListener(t *testing.T) {
	mockConsumer := &MockConsumer{}
	listener, err := NewMessageListener("test_index", mockConsumer, MockParseFunc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if listener.IndexName != "test_index" {
		t.Errorf("expected index name to be 'test_index', got %s", listener.IndexName)
	}
}

func TestStartStopMessageListener(t *testing.T) {
	mockConsumer := &MockConsumer{}
	mockey.Mock(elastic.ExistIndex).Return(true, nil).Build()
	defer mockey.UnPatchAll()

	listener, err := NewMessageListener("test_index", mockConsumer, MockParseFunc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = listener.Start()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	time.Sleep(1 * time.Second) // Allow some fetch/sync cycles to execute

	err = listener.Stop()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMessageFetch(t *testing.T) {
	mockConsumer := &MockConsumer{}

	listener, err := NewMessageListener("test_index", mockConsumer, MockParseFunc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	listener.fetch()
	if len(listener.messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(listener.messages))
	}

	var ok bool
	for _, message := range MessageSet {
		if string(message.GetValue()) == string(listener.messages[0].GetValue()) {
			ok = true
			break
		}
	}

	if !ok {
		t.Errorf("expected message to be in MessageSet, got %s", string(listener.messages[0].GetValue()))
	}
}

func TestMessageSync(t *testing.T) {
	const messageCount = 5

	mockConsumer := &MockConsumer{}
	mockey.Mock(elastic.IndexDocuments).To(func(ctx context.Context, indexName string, documents []elastic.IDocument) error {
		if indexName != "test_index" {
			t.Errorf("expected index name to be 'test_index', got %s", indexName)
		}

		if len(documents) != messageCount {
			t.Errorf("expected document count to be 5, got %d", len(documents))
		}

		return nil
	}).Build()

	defer mockey.UnPatchAll()

	listener, err := NewMessageListener("test_index", mockConsumer, MockParseFunc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 0; i < messageCount; i++ {
		listener.messages = append(listener.messages, MessageSet[i])
		listener.messages = append(listener.messages, MessageSet[i]) // duplicated message.
	}

	listener.sync()

	if len(listener.messages) != 0 {
		t.Errorf("expected messages to be empty, got %d", len(listener.messages))
	}
}