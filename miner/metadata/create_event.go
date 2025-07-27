package metadata

import (
	"fmt"

	"github.com/google/uuid"
	api "github.com/kinneko-de/sample-eventual-consistency-transaction-log-tailing-mongodb/golang/store_file/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateFileStoredEvent(change bson.M) (*api.FileStored, error) {
	// remarks: this works only for our scenario, where we have only update events with fullDocument configured
	if fullDoc, ok := change["fullDocument"].(bson.M); ok {
		fmt.Printf("Full document: %v\n", fullDoc)
		fileStoredEvent := &api.FileStored{}

		fileIdBin, ok := fullDoc["FileId"].(primitive.Binary)
		if !ok || fileIdBin.Subtype != 4 || len(fileIdBin.Data) != 16 {
			return nil, fmt.Errorf("FileId missing or not a valid UUID binary")
		}
		u, err := uuid.FromBytes(fileIdBin.Data)
		if err != nil {
			return nil, fmt.Errorf("FileId bytes could not be parsed as UUID: %w", err)
		}
		fileStoredEvent.SetFileId(u.String())

		createdAt, ok := fullDoc["CreatedAt"].(primitive.DateTime)
		if !ok {
			return nil, fmt.Errorf("CreatedAt missing or not a primitive.DateTime")
		}
		fileStoredEvent.SetCreatedAt(timestamppb.New(createdAt.Time()))

		storedAt, ok := fullDoc["StoredAt"].(primitive.DateTime)
		if !ok {
			return nil, fmt.Errorf("StoredAt missing or not a primitive.DateTime")
		}
		fileStoredEvent.SetStoredAt(timestamppb.New(storedAt.Time()))

		size, ok := fullDoc["Size"].(int64)
		if !ok {
			return nil, fmt.Errorf("Size missing or not an int64")
		}
		fileStoredEvent.SetSize(size)

		mediaType, ok := fullDoc["MediaType"].(string)
		if !ok {
			return nil, fmt.Errorf("MediaType missing or not a string")
		}
		fileStoredEvent.SetMediaType(mediaType)

		extension, ok := fullDoc["Extension"].(string)
		if !ok {
			return nil, fmt.Errorf("Extension missing or not a string")
		}
		fileStoredEvent.SetExtension(extension)

		fmt.Printf("FileStored event: %+v\n", fileStoredEvent)

		return fileStoredEvent, nil
	}

	return nil, fmt.Errorf("Scenario not supported, only events with fullDocument configured are currently supported")

}
