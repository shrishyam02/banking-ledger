package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type kafkaConsumer struct {
	reader *kafka.Reader
}

type KafkaConsumer interface {
	Consume(ctx context.Context, topic string, groupID string, handler func(kafka.Message) error) error
}

func NewKafkaConsumer(brokers []string, groupID string, groupTopics []string) KafkaConsumer {
	return &kafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			GroupID:     groupID,
			GroupTopics: groupTopics,
		}),
	}
}

func (kc *kafkaConsumer) Consume(ctx context.Context, topic string, groupID string, handler func(kafka.Message) error) error {
	kc.reader.SetOffset(kafka.FirstOffset)
	for {
		msg, err := kc.reader.FetchMessage(ctx)
		if err != nil {
			return err
		}
		if err := handler(msg); err != nil {
			log.Printf("Failed to handle message: %v", err)
			continue
		}
		if err := kc.reader.CommitMessages(ctx, msg); err != nil {
			return err
		}
	}
}
