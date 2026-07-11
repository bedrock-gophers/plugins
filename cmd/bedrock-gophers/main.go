package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bedrock-gophers/plugins/framework"
)

func main() {
	configPath := flag.String("config", "server.toml", "path to server configuration")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := framework.RunFile(ctx, *configPath, slog.Default()); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
