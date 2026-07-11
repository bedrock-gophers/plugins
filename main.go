package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/bedrock-gophers/plugins/framework"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := framework.RunFile(ctx, "server.toml", nil); err != nil {
		log.Fatal(err)
	}
}

