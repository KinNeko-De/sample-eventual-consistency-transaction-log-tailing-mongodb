package clean

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type IncompleteMetadata struct {
	FileId    uuid.UUID
	CreatedAt time.Time
}

func CleanWhatWasLeftBehind(ctx context.Context) error {
	fmt.Println("Cleaning...")

	err := InitializeMongoClient(ctx)
	if err != nil {
		return fmt.Errorf("Error initializing MongoDB client: %w", err)
	}
	defer DisconnectMongoClient()

	_, err = FetchIncompleteMetadata(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching incomplete metadata: %w", err)
	}

	fmt.Println("Cleaned")
	return nil
}
