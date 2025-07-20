package file

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StoreFileId(fileId uuid.UUID) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), nil
}

func StoreFileMetadata(objectId primitive.ObjectID, size uint64) error {
	return nil
}
