package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type kafkaProducer struct {
	writer *kafka.Writer
}

type KafkaProducer interface {
	Produce(ctx context.Context, topic string, message kafka.Message) error
}

func NewKafkaProducer(brokers []string) KafkaProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (kp *kafkaProducer) Produce(ctx context.Context, topic string, message kafka.Message) error {
	message.Topic = topic
	return kp.writer.WriteMessages(ctx, message)
}
