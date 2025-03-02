package kafka

import (
	"log"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
	"github.com/shrishyam02/banking-ledger/common/logger"
)

type kafkaTopic struct{}

type KafkaTopic interface {
	CreateKafkaTopic(brokers, topic string) error
}

func CreateKafkaTopic(brokers, topic string) error {
	conn, err := kafka.Dial("tcp", brokers)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	// Check if the topic already exists
	partitions, err := controllerConn.ReadPartitions(topic)
	logger.Log.Info().Msgf("Partitions: %v", partitions)
	if err == nil && len(partitions) > 0 {
		log.Printf("Topic %s already exists\n", topic)
		return nil
	}

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	return controllerConn.CreateTopics(topicConfigs...)
}
