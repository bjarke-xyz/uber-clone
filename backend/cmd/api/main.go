package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/bjarke-xyz/uber-clone-backend/internal/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	err := cmd.APICmd(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
