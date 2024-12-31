package messages

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type KafkaMessage kafka.Message

type KafkaConsumer struct {
	reader *kafka.Reader
}

func (m KafkaMessage) GetValue() []byte {
	return m.Value
}

func (kc *KafkaConsumer) FetchMessage(ctx context.Context) (IMessage, error) {
	msg, err := kc.reader.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}

	return KafkaMessage(msg), nil
}

func (kc *KafkaConsumer) CommitMessages(ctx context.Context, messages ...IMessage) error {
	msgs := make([]kafka.Message, 0, len(messages))
	for _, msg := range messages {
		msgs = append(msgs, kafka.Message(msg.(KafkaMessage)))
	}

	return kc.reader.CommitMessages(ctx, msgs...)
}

func (kc *KafkaConsumer) Close() error {
	return kc.reader.Close()
}

func NewKafkaConsumer(brokers []string, topic, groupID string, startOffset int64) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       topic,
			GroupID:     groupID,
			StartOffset: startOffset,
		}),
	}
}
