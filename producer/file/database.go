package file

import (
	"fmt"
	"math/rand/v2"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrorProbabilityFileId   float64 = 0.01
	ErrorProbabilityMetadata float64 = 0.1
)

func StoreFileId(fileId uuid.UUID) (primitive.ObjectID, error) {
	if rand.Float64() < ErrorProbabilityFileId {
		return primitive.NilObjectID, fmt.Errorf("Failed to write to database")
	}

	return primitive.NewObjectID(), nil
}

func StoreFileMetadata(objectId primitive.ObjectID, size uint64, mediaType string) error {
	if rand.Float64() < ErrorProbabilityMetadata {
		return fmt.Errorf("Failed to write to database")
	}

	return nil
}
