package bootstrap

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	chatrepository "dan-ai/internal/chat/repository"
	"dan-ai/internal/memory/extractor"
	"dan-ai/internal/memory/processor"
	memoryrepo "dan-ai/internal/memory/repository"
	memoryservice "dan-ai/internal/memory/service"
	"dan-ai/pkg/config"
	"dan-ai/pkg/kafka"
	"dan-ai/pkg/postgres"
)

const GroupID = "dan-ai-memory-worker"

func RunMemoryWorker() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := postgres.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	chatRepo := chatrepository.NewPostgresRepository(db)
	memoryRepo := memoryrepo.NewPostgresRepository(db)
	memorySvc := memoryservice.NewService(memoryRepo)
	extractor := extractor.NewExtractor()
	processor := processor.NewProcessor(chatRepo, memorySvc, extractor)

	consumer := kafka.NewConsumer(cfg, "portfolio.events", GroupID)
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("error closing consumer: %v", err)
		}
	}()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("shutting down memory worker...")
		cancel()
	}()

	log.Println("starting memory worker, waiting for chat.completed events...")
	if err := consumer.Consume(ctx, processor.ProcessEvent); err != nil && err != context.Canceled {
		log.Fatalf("consumer error: %v", err)
	}
}
