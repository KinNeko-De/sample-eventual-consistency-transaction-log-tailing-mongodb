package file

import (
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"path"

	"github.com/google/uuid"
)

var (
	ErrorProbabilityFolderCreate         float64 = 0.01
	ErrorProbabilityFileCreate           float64 = 0.01
	ErrorProbabilityFileClose            float64 = 0.01
	ErrorProbabilityFileWriteByte        float64 = 0.02
	ErrorProbabilityUserPrMetworkAborted float64 = 0.05
)

const (
	StoragePath = "storage"
)

func StoreFileBytes(fileId uuid.UUID) (uint64, string, error) {
	writeCloser, err := CreateFileAndFolder(fileId)
	if err != nil {
		return 0, "", err
	}

	fileSize, mediaType, err := WriteChunk(writeCloser)
	if err != nil {
		return 0, "", err
	}

	err = CloseFile(writeCloser)
	if err != nil {
		return 0, "", err
	}

	return fileSize, mediaType, nil
}

func CreateFileAndFolder(fileId uuid.UUID) (io.WriteCloser, error) {
	fileFolder := path.Join(StoragePath, fileId.String())
	fileLocation := path.Join(fileFolder, fileId.String())

	err := CreateFolder(fileFolder)
	if err != nil {
		return nil, err
	}

	wc, err := CreateFile(fileLocation)
	return wc, err
}

func CreateFolder(fileFolder string) error {
	if rand.Float64() < ErrorProbabilityFolderCreate {
		return fmt.Errorf("Failed to create folder")
	}
	err := os.MkdirAll(fileFolder, os.ModePerm)
	return err
}

func CreateFile(fileLocation string) (io.WriteCloser, error) {
	if rand.Float64() < ErrorProbabilityFileCreate {
		return nil, fmt.Errorf("Failed to create file")
	}

	writer, err := os.Create(fileLocation)
	return writer, err
}

// WriteChunk simulates writing a file in randomly sized chunks, introducing random errors to mimic
// user aborts, network or file write failures. It returns the total file size written, the media type of the file,
// and an error if any failure occurs during the process.
//
// Returns:
//   - uint64: The total size of the file written in bytes.
//   - string: The media type of the file
//   - error: An error if the stream is aborted or a file write fails; otherwise, nil.
func WriteChunk(fileWriter io.WriteCloser) (uint64, string, error) {
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

		randomText := GenerateRandomText(chunkSize)
		_, err := fileWriter.Write(randomText)
		if err != nil {
			return 0, "", fmt.Errorf("Failed to write file bytes: %w", err)
		}
	}

	return fileSize, "text/plain", nil
}

func GenerateRandomText(chunkSize uint64) []byte {
	letters := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 ")
	randomText := make([]byte, chunkSize)
	for i := range randomText {
		randomText[i] = letters[rand.IntN(len(letters))]
	}
	return randomText
}

func CloseFile(fileCloser io.WriteCloser) error {
	if rand.Float64() < ErrorProbabilityFileClose {
		return fmt.Errorf("Failed to write to database")
	}

	closeErr := fileCloser.Close()
	return closeErr
}
