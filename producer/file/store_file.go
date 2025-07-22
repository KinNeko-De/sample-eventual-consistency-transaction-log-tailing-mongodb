package file

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
)

var (
	minDelay = 2
	maxDelay = 5
)

func SimulateStoreFile(ctx context.Context) error {
	fmt.Println("Storing File...")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, storing File stopped")
			return ctx.Err()
		case <-time.After(CreateJitteredDelay()):
			err := StoreFile(ctx)
			if err != nil && (os.IsTimeout(err) || err == context.Canceled || err == context.DeadlineExceeded) {
				return err
			}
			if err != nil {
				fmt.Printf("Error storing file: %v", err)
			} else {
				fmt.Println("File stored successfully")
			}
		}
	}
}

func CreateJitteredDelay() time.Duration {
	jitter := time.Duration(rand.Intn(maxDelay-minDelay)+minDelay) * time.Second
	fmt.Printf("Jittered delay: %v\n", jitter)
	return jitter
}

func StoreFile(ctx context.Context) error {
	fileId := uuid.New()
	objectId, err := StoreFileId(fileId)
	if err != nil {
		return fmt.Errorf("Error storing file ID: %w for %v\n", err, fileId)
	}
	size, mediaType, err := StoreFileBytes(fileId)
	if err != nil {
		return fmt.Errorf("Error storing file bytes: %w for %v\n", err, fileId)
	}
	err = StoreFileMetadata(objectId, size, mediaType)
	if err != nil {
		return fmt.Errorf("Error storing file metadata: %w for %v\n", err, fileId)
	}
	return nil
}
