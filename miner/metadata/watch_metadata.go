package metadata

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
)

func MiningFileMetadata(ctx context.Context) error {
	fmt.Println("Mining file metadata...")

	if err := initializeMongoClient(ctx); err != nil {
		return err
	}
	defer disconnectMongoClient()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, mining file metadata stopped")
			return ctx.Err()
		default:
			if err := WatchChangeStream(ctx); err != nil {
				if ctx.Err() != nil && (ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded) {
					return ctx.Err()
				}

				return fmt.Errorf("failed to watch file metadata: %w", err)
			}
		}
	}
}

func WatchChangeStream(ctx context.Context) error {
	fmt.Println("Watching change stream for file metadata...")

	collection := client.Database("store_file").Collection("file")
	changeStreamOptions := options.ChangeStream().SetFullDocument(options.Required)
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "operationType", Value: "update"},
			{Key: "updateDescription.updatedFields.StoredAt", Value: bson.D{{Key: "$exists", Value: true}}},
		}}},
	}
	changeStream, err := collection.Watch(ctx, pipeline, changeStreamOptions)
	if err != nil {
		return fmt.Errorf("failed to watch change stream: %w", err)
	}
	defer changeStream.Close(ctx)

	for changeStream.Next(ctx) {
		var change bson.M
		if err := changeStream.Decode(&change); err != nil {
			return fmt.Errorf("failed to decode change stream event: %w", err)
		}
		fmt.Printf("Change detected: %v\n", change)
	}

	if err := changeStream.Err(); err != nil {
		if err == context.Canceled {
			return ctx.Err()
		}
		return fmt.Errorf("error in change stream: %w", err)
	}

	return nil
}

func initializeMongoClient(ctx context.Context) error {
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

func disconnectMongoClient() {
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
