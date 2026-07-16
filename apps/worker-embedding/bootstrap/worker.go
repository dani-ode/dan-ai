package bootstrap

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	aiClient "dan-ai/internal/ai/client"
	"dan-ai/internal/ai/provider"
	"dan-ai/internal/knowledge/chunk"
	"dan-ai/internal/knowledge/processor"
	"dan-ai/internal/knowledge/repository"
	"dan-ai/pkg/config"
	"dan-ai/pkg/kafka"
	"dan-ai/pkg/milvus"
	"dan-ai/pkg/postgres"
)

const GroupID = "portfolio-embedding-worker"

func RunEmbeddingWorker() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Postgres
	db, err := postgres.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// Initialize Milvus
	milvusCtx, milvusCancel := context.WithTimeout(ctx, 5*time.Second)
	mClient, err := milvus.NewClient(milvusCtx, cfg)
	milvusCancel()
	if err != nil {
		log.Fatalf("failed to connect to milvus: %v", err)
	}
	if err := mClient.InitCollection(ctx); err != nil {
		log.Fatalf("failed to init milvus collection: %v", err)
	}

	// Initialize AI
	genaiClient, err := aiClient.NewClient(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to init ai client: %v", err)
	}
	aiProvider := provider.NewGeminiProvider(genaiClient)

	// Initialize chunk builder
	chunkBuilder := chunk.NewAIBuilder(aiProvider)

	// Initialize Knowledge Processor
	repo := repository.NewPostgresKnowledgeRepository(db)
	proc := processor.NewProcessor(repo, aiProvider, mClient, chunkBuilder)

	// Initialize Kafka Consumer
	consumer := kafka.NewConsumer(cfg, "portfolio.knowledge", GroupID)
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("error closing consumer: %v", err)
		}
	}()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("shutting down embedding worker...")
		cancel()
	}()

	// Start consuming
	log.Println("starting embedding worker, waiting for events...")
	err = consumer.Consume(ctx, proc.ProcessEvent)
	if err != nil && err != context.Canceled {
		log.Fatalf("consumer error: %v", err)
	}
}
