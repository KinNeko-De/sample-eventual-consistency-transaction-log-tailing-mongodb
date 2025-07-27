package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/KinNeko-De/sample-eventual-consistency-transaction-log-tailing-mongodb/consumer/metadata"
)

func main() {
	fmt.Println("Starting consumer...")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := metadata.ConsumingFileStored(ctx)
		if err != nil && !os.IsTimeout(err) && err != context.Canceled && err != context.DeadlineExceeded {
			fmt.Printf("Error consuming file metadata: %v\n", err)
		}
	}()

	wg.Wait()

	fmt.Println("Shutting down consumer...")
}
