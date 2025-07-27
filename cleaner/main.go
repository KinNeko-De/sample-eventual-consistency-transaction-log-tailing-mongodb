package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/KinNeko-De/sample-eventual-consistency-transaction-log-tailing-mongodb/cleaner/clean"
)

func main() {
	fmt.Println("Starting cleaner...")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	err := clean.CleanWhatWasLeftBehind(ctx)
	if err != nil {
		fmt.Printf("Error during cleaning: %v\n", err)
	}

	fmt.Println("Shutting down cleaner...")
}
