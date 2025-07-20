package file

import (
	"context"
	"fmt"
	"math/rand"
	"time"
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
			fmt.Println("File stored")
		}
	}
}

func CreateJitteredDelay() time.Duration {
	jitter := time.Duration(rand.Intn(maxDelay-minDelay)+minDelay) * time.Second
	fmt.Printf("Jittered delay: %v\n", jitter)
	return jitter
}
