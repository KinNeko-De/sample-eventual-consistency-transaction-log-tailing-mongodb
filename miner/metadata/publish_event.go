package metadata

import (
	"fmt"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

var (
	brokers  = []string{"localhost:9095"}
	topic    = "file-stored"
	producer sarama.SyncProducer
)

func PublishEvent(event proto.Message) error {
	fmt.Println("Publishing event...")

	msgBytes, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf message: %w", err)
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msgBytes),
	}

	partition, offset, err := producer.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("failed to publish event to Kafka: %w", err)
	}

	fmt.Printf("Event published to partition %d at offset %d\n", partition, offset)
	return nil
}

func CreateEventProducer() error {
	if producer == nil {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true

		var err error
		producer, err = sarama.NewSyncProducer(brokers, config)
		if err != nil {
			return fmt.Errorf("failed to create producer: %w", err)
		}

		fmt.Println("Kafka producer created")
	}

	return nil
}

func CloseEventProducer() error {
	if producer != nil {
		if err := producer.Close(); err != nil {
			return fmt.Errorf("Failed to close producer: %w\n", err)
		}
	}

	return nil
}
