package clean

import (
	"context"
	"fmt"
)

func CleanWhatWasLeftBehind(ctx context.Context) error {
	fmt.Println("Cleaning...")

	err := InitializeMongoClient(ctx)
	if err != nil {
		return fmt.Errorf("Error initializing MongoDB client: %w", err)
	}
	defer DisconnectMongoClient()

	fmt.Println("Cleaned")
	return nil
}
