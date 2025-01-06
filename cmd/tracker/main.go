package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gocoinspot/internal/collector"
	"gocoinspot/internal/config"
	"gocoinspot/internal/storage"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fs, err := storage.NewFileSystem(cfg.OutputDir)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	collector := collector.NewCollector(cfg, fs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		cancel()
	}()

	if err := collector.Start(ctx); err != nil {
		log.Fatalf("Collector stopped: %v", err)
	}
}
