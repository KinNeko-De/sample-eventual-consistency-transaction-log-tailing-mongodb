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

	files, err := FetchIncompleteMetadata(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching incomplete metadata: %w", err)
	}

	for _, file := range files {
		err := CleanFileBytes(file)
		if err != nil {
			return fmt.Errorf("Error cleaning file bytes for FileId %s: %w", file.FileId.String(), err)
		}
		err = CleanFileMetadata(ctx, file)
		if err != nil {
			return fmt.Errorf("Error deleting incomplete metadata for FileId %s: %w", file.FileId.String(), err)
		}
	}

	fmt.Println("Cleaned")
	return nil
}
