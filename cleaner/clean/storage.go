package clean

import (
	"fmt"
	"os"
	"path"
)

const (
	StoragePath = "../producer/storage"
)

func CleanFileBytes(file IncompleteMetadata) error {
	fileFolder := path.Join(StoragePath, file.FileId.String())

	err := os.RemoveAll(fileFolder)
	if err != nil {
		return fmt.Errorf("failed to remove folder %s: %w", fileFolder, err)
	}

	fmt.Printf("Cleaned up file bytes for FileId: %s\n", file.FileId.String())
	return nil
}
