package clean

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client    *mongo.Client
	OlderThan time.Duration = time.Hour
)

const Limit int64 = 1000

func FetchIncompleteMetadata(ctx context.Context) ([]IncompleteMetadata, error) {
	collection := client.Database("store_file").Collection("file")

	cutoff := time.Now().UTC().Add(-OlderThan)
	filter := map[string]any{
		"StoredAt":  map[string]any{"$exists": false},
		"CreatedAt": map[string]any{"$lt": cutoff},
	}

	findOpts := options.Find().SetLimit(Limit)
	cursor, err := collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to query incomplete metadata: %w", err)
	}
	defer cursor.Close(ctx)

	var results []IncompleteMetadata
	for cursor.Next(ctx) {
		var rawDoc bson.Raw
		if err := cursor.Decode(&rawDoc); err != nil {
			return nil, fmt.Errorf("failed to decode raw document: %w", err)
		}
		document, err := UnmarshalBSON(rawDoc)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal IncompleteMetadata: %w", err)
		}

		results = append(results, document)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	fmt.Printf("Found %d incomplete metadata entries older than %s\n", len(results), OlderThan)
	for _, entry := range results {
		fmt.Printf("FileId: %s, CreatedAt: %s (UTC)\n", entry.FileId, entry.CreatedAt.UTC().Format(time.RFC3339))
	}
	return results, nil
}

func CleanFileMetadata(ctx context.Context, file IncompleteMetadata) error {
	collection := client.Database("store_file").Collection("file")

	filter := bson.M{"FileId": primitive.Binary{Subtype: 4, Data: file.FileId[:]}}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete incomplete metadata for FileId %s: %w", file.FileId.String(), err)
	}

	fmt.Printf("Deleted incomplete metadata for FileId %s\n", file.FileId.String())
	return nil
}

func UnmarshalBSON(data []byte) (IncompleteMetadata, error) {
	raw := bson.Raw(data)

	fileIdVal := raw.Lookup("FileId")
	_, fieldData := fileIdVal.Binary()
	id, err := uuid.FromBytes(fieldData)
	if err != nil {
		return IncompleteMetadata{}, fmt.Errorf("FileId is not a valid UUID: %w", err)
	}

	createdAtVal := raw.Lookup("CreatedAt")
	if createdAtVal.Type != bson.TypeDateTime {
		return IncompleteMetadata{}, fmt.Errorf("CreatedAt is not a datetime type")
	}
	createdAt := createdAtVal.Time()

	return IncompleteMetadata{
		FileId:    id,
		CreatedAt: createdAt,
	}, nil
}

func InitializeMongoClient(ctx context.Context) error {
	if client == nil {
		var err error
		clientOptions := options.Client().
			ApplyURI("mongodb://localhost:27017/?replicaSet=rs0")
		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			return fmt.Errorf("failed to connect to MongoDB: %v", err)
		}
		fmt.Println("MongoDB client initialized")

		if err := client.Ping(ctx, nil); err != nil {
			return fmt.Errorf("failed to ping MongoDB: %v", err)
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
