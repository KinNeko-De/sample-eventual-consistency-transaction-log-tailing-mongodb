package metadata

import (
	"context"
	"fmt"
)

func ConsumingFileStored(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, consuming file metadata stopped")
			return ctx.Err()
		default:
			fmt.Println("Waiting for file metadata..")
		}
	}
}
