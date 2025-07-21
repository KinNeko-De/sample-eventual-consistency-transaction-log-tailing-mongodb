package file

import (
	"fmt"
	"math/rand/v2"

	"github.com/google/uuid"
)

var (
	ErrorProbabilityFolderCreate         float64 = 0.01
	ErrorProbabilityFileCreate           float64 = 0.01
	ErrorProbabilityFileClose            float64 = 0.01
	ErrorProbabilityFileWriteByte        float64 = 0.02
	ErrorProbabilityUserPrMetworkAborted float64 = 0.05
)

func StoreFileBytes(fileId uuid.UUID) (uint64, string, error) {
	err := CreateFileAndFolder(fileId)
	if err != nil {
		return 0, "", err
	}

	fileSize, mediaType, err := WriteChunk()
	if err != nil {
		return 0, "", err
	}

	err = CloseFile()
	if err != nil {
		return 0, "", err
	}

	return fileSize, mediaType, nil
}

func CreateFileAndFolder(fileId uuid.UUID) error {
	if rand.Float64() < ErrorProbabilityFolderCreate {
		return fmt.Errorf("Failed to create folder")
	}
	if rand.Float64() < ErrorProbabilityFileCreate {
		return fmt.Errorf("Failed to create file")
	}
	return nil
}

// WriteChunk simulates writing a file in randomly sized chunks, introducing random errors to mimic
// user aborts, network or file write failures. It returns the total file size written, the media type of the file,
// and an error if any failure occurs during the process.
//
// Returns:
//   - uint64: The total size of the file written in bytes.
//   - string: The media type of the file
//   - error: An error if the stream is aborted or a file write fails; otherwise, nil.
func WriteChunk() (uint64, string, error) {
	var fileSize uint64 = 0
	numberOfChunks := rand.IntN(4) + 2
	for i := 0; i < numberOfChunks; i++ {
		if rand.Float64() < ErrorProbabilityUserPrMetworkAborted {
			return 0, "", fmt.Errorf("Stream was aborted")
		}
		chunkSize := rand.Uint64N(1024) + 512 // Random chunk size between 512 bytes and 1536 bytes
		fileSize += chunkSize
		if rand.Float64() < ErrorProbabilityFileWriteByte {
			return 0, "", fmt.Errorf("Failed to write file bytes")
		}
	}

	return fileSize, "text/plain", nil
}

func CloseFile() error {
	if rand.Float64() < ErrorProbabilityFileClose {
		return fmt.Errorf("Failed to write to database")
	}

	return nil
}
