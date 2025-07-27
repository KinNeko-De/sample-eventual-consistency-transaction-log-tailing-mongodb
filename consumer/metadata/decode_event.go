package metadata

import (
	"fmt"

	"github.com/IBM/sarama"
	api "github.com/kinneko-de/sample-eventual-consistency-transaction-log-tailing-mongodb/golang/store_file/v1"
	"google.golang.org/protobuf/proto"
)

func DecodeMessage(message *sarama.ConsumerMessage) error {
	fileStored := &api.FileStored{}
	err := proto.Unmarshal(message.Value, fileStored)
	if err != nil {
		return fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}
	fmt.Printf("Consumed FileStored: %+v\n", fileStored)

	return nil
}
