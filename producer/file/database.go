package file

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrorProbabilityFileId   float64 = 0.01
	ErrorProbabilityMetadata float64 = 0.1
	client                   *mongo.Client
)

func StoreFileId(ctx context.Context, fileId uuid.UUID) (primitive.ObjectID, error) {
	if rand.Float64() < ErrorProbabilityFileId {
		return primitive.NilObjectID, fmt.Errorf("Failed to write to database")
	}

	collection := client.Database("store_file").Collection("file")

	objectId := primitive.NewObjectID()

	document := bson.M{
		"_id":       objectId,
		"FileId":    primitive.Binary{Subtype: 4, Data: fileId[:]}, // store UUID as BSON binary subtype 4
		"CreatedAt": time.Now().UTC(),
	}

	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert file id: %w", err)
	}

	return objectId, nil
}

func StoreFileMetadata(objectId primitive.ObjectID, size uint64, mediaType string) error {
	if rand.Float64() < ErrorProbabilityMetadata {
		return fmt.Errorf("Failed to write to database")
	}

	return nil
}

func InitializeMongoClient(ctx context.Context) error {
	if client == nil {
		var err error
		clientOptions := options.Client().
			ApplyURI("mongodb://localhost:27017/?replicaSet=rs0")

		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			return fmt.Errorf("failed to connect to MongoDB: %w", err)
		}

		fmt.Println("MongoDB client initialized")

		if err := client.Ping(ctx, nil); err != nil {
			return fmt.Errorf("failed to ping MongoDB: %w", err)
		}

		fmt.Println("MongoDB ping successful")
	}

	return nil
}

func DisconnectMongoClient() {
	if client != nil {
		// Create a new context with timeout for disconnect operation, the application context might be cancelled
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		if err := client.Disconnect(disconnectCtx); err != nil {
			fmt.Printf("Failed to disconnect MongoDB client: %v\n", err)
		} else {
			fmt.Println("MongoDB client disconnected")
		}
	}
}
