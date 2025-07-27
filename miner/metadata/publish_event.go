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

	fmt.Println("Event publshed")

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
