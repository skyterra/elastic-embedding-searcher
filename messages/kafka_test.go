package messages

import (
	"context"
	"github.com/segmentio/kafka-go"
	"sync"
	"testing"
)

const (
	groupID = "kafka-test"
)

func TestKafkaConsumer(t *testing.T) {
	consumer := NewKafkaConsumer([]string{"localhost:9092"}, "quotes-test", groupID, kafka.LastOffset)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		message, err := consumer.FetchMessage(context.Background())
		if err != nil {
			t.Errorf("fail to fetch message. err:%s", err.Error())
		}

		if string(message.GetValue()) != string(MessageSet[0].GetValue()) {
			t.Errorf("expected message to be '%s', got %s", MessageSet[0].GetValue(), message.GetValue())
		}

		consumer.CommitMessages(context.Background(), message)
		t.Log(string(message.GetValue()))
	}()

	writer := &kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"),
		Topic: "quotes-test",
	}

	err := writer.WriteMessages(context.Background(), kafka.Message{
		Value: MessageSet[0].GetValue(),
	})
	if err != nil {
		t.Errorf("fail to write message. err:%s", err.Error())
	}

	wg.Wait()
}
