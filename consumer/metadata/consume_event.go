package metadata

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

var (
	brokers       = []string{"localhost:9095"}
	topic         = "file-stored"
	groupID       = "file-stored-group"
	consumerGroup sarama.ConsumerGroup
)

type fileStoredHandler struct{}

func (h *fileStoredHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *fileStoredHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *fileStoredHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := DecodeMessage(message); err != nil {
			fmt.Printf("Error decoding message: %v\n", err)
		}
		sess.MarkMessage(message, "")
	}
	return nil
}

func ConsumingFileStored(ctx context.Context) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	var err error
	consumerGroup, err = sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	defer consumerGroup.Close()

	handler := &fileStoredHandler{}

	fmt.Println("Waiting file metadata events...")
	for {
		if err := consumerGroup.Consume(ctx, []string{topic}, handler); err != nil {
			return fmt.Errorf("error from consumer: %w", err)
		}
		if ctx.Err() != nil {
			fmt.Println("Context cancelled, mining file metadata stopped")
			return ctx.Err()
		}
	}
}
