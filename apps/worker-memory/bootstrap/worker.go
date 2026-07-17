package bootstrap

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	aiClient "dan-ai/internal/ai/client"
	aiprovider "dan-ai/internal/ai/provider"
	chatrepository "dan-ai/internal/chat/repository"
	embeddingrepo "dan-ai/internal/embedding/repository"
	"dan-ai/internal/memory/extractor"
	"dan-ai/internal/memory/processor"
	memoryrepo "dan-ai/internal/memory/repository"
	memoryservice "dan-ai/internal/memory/service"
	promptrepo "dan-ai/internal/prompt/repository"
	"dan-ai/pkg/config"
	"dan-ai/pkg/kafka"
	"dan-ai/pkg/milvus"
	"dan-ai/pkg/postgres"
)

const GroupID = "dan-memory-worker"

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

	// Initialize embedding repo and active profile config
	embeddingRepo := embeddingrepo.NewPostgresRepository(db)
	profile, err := embeddingRepo.GetActiveProfile(ctx)
	if err != nil {
		log.Fatalf("failed to get active embedding profile: %v", err)
	}

	// Initialize Milvus Client
	milvusCtx, milvusCancel := context.WithTimeout(ctx, 5*time.Second)
	mClient, err := milvus.NewClient(milvusCtx, cfg)
	milvusCancel()
	if err != nil {
		log.Fatalf("failed to connect to milvus: %v", err)
	}
	if err := mClient.InitCollection(ctx, profile.KnowledgeCollection, profile.VisitorCollection, profile.Dimension, profile.MetricType); err != nil {
		log.Fatalf("failed to init milvus collection: %v", err)
	}

	// Initialize AI provider registry
	aiRegistry := aiprovider.NewRegistry()

	genaiClient, err := aiClient.NewClient(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialize Gemini AI client: %v", err)
	}
	aiRegistry.Register("gemini", aiprovider.NewGeminiProvider(genaiClient))

	if cfg.AI.OpenAIAPIKey != "" {
		aiRegistry.Register("openai", aiprovider.NewOpenAIProvider(cfg.AI.OpenAIAPIKey))
		log.Println("OpenAI provider registered")
	}

	chatRepo := chatrepository.NewPostgresRepository(db)
	memoryRepo := memoryrepo.NewPostgresRepository(db)
	promptRepo := promptrepo.NewPostgresRepository(db)
	memorySvc := memoryservice.NewService(memoryRepo, mClient, aiRegistry, promptRepo, embeddingRepo)
	memoryExtractor := extractor.NewExtractor(aiRegistry, promptRepo)
	processor := processor.NewProcessor(chatRepo, memorySvc, memoryExtractor, promptRepo)

	consumer := kafka.NewConsumer(cfg, "dan.events", GroupID)
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
