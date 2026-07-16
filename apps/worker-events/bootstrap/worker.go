package bootstrap

import (
	"context"
	"log"
	"os"
	"os/signal"
	"dan-ai/internal/outbox/publisher"
	"dan-ai/internal/outbox/repository"
	"dan-ai/pkg/config"
	"dan-ai/pkg/kafka"
	"dan-ai/pkg/postgres"
	"syscall"
	"time"
)

func RunEventWorker() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize Postgres
	db, err := postgres.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Initialize Kafka Producer
	producer := kafka.NewProducer(cfg)
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("error closing producer: %v", err)
		}
	}()

	// Initialize Outbox Repository and Publisher
	outboxRepo := repository.NewPostgresRepository(db)
	outboxPublisher := publisher.NewPublisher(outboxRepo, producer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("shutting down event worker...")
		cancel()
	}()

	// Start publisher loop
	log.Println("starting event worker...")
	outboxPublisher.Start(ctx, 5*time.Second) // Poll every 5 seconds
}
