package metadata

import (
	"context"
	"fmt"
)

func MiningFileMetadata(ctx context.Context) error {
	fmt.Println("Mining file metadata...")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, mining file metadata stopped")
			return ctx.Err()
		default:
			fmt.Println("Watching for file metadata..")
		}
	}
}
