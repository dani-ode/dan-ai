package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"dan-ai/internal/ai/provider"
	embeddingrepo "dan-ai/internal/embedding/repository"
	"dan-ai/internal/knowledge/chunk"
	"dan-ai/internal/knowledge/entity"
	"dan-ai/internal/knowledge/repository"
	promptrepo "dan-ai/internal/prompt/repository"
	"dan-ai/pkg/kafka"
	"dan-ai/pkg/milvus"
	"dan-ai/pkg/ulid"
)

const defaultChunkModel = "gemini-3.1-flash-lite"

type Processor struct {
	repo          repository.KnowledgeRepository
	aiRegistry    *provider.Registry
	milvusClient  *milvus.Client
	chunkBuilder  chunk.Builder
	promptRepo    promptrepo.Repository
	embeddingRepo embeddingrepo.Repository
}

func NewProcessor(
	repo repository.KnowledgeRepository,
	aiRegistry *provider.Registry,
	milvusClient *milvus.Client,
	chunkBuilder chunk.Builder,
	promptRepo promptrepo.Repository,
	embeddingRepo embeddingrepo.Repository,
) *Processor {
	return &Processor{
		repo:          repo,
		aiRegistry:    aiRegistry,
		milvusClient:  milvusClient,
		chunkBuilder:  chunkBuilder,
		promptRepo:    promptRepo,
		embeddingRepo: embeddingRepo,
	}
}

type knowledgeEventPayload struct {
	SourceType string `json:"source_type"`
	SourceID   string `json:"source_id"`
	PromptID   string `json:"prompt_id"`
}

func (p *Processor) ProcessEvent(ctx context.Context, event kafka.Event) error {
	log.Printf("processing knowledge event for %s %s", event.Aggregate, event.AggregateID)

	// Get active embedding profile
	profile, err := p.embeddingRepo.GetActiveProfile(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active embedding profile: %w", err)
	}

	// Parse payload to extract prompt_id
	var payload knowledgeEventPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		log.Printf("warning: failed to parse knowledge event payload: %v", err)
	}

	// Resolve model name, provider, and system instruction from prompt
	modelName := defaultChunkModel
	chunkProviderName := "gemini" // default provider for chunk generation
	systemInstruction := ""
	if payload.PromptID != "" {
		prompt, err := p.promptRepo.Get(ctx, payload.PromptID)
		if err == nil {
			if prompt.AIModel.Name != "" {
				modelName = prompt.AIModel.Name
			}
			if prompt.AIModel.Provider != "" {
				chunkProviderName = prompt.AIModel.Provider
			}
			systemInstruction = prompt.SystemPrompt
		}
	}

	// Fetch document by source
	doc, err := p.repo.GetDocumentBySource(ctx, event.Aggregate, event.AggregateID)
	if err != nil {
		return fmt.Errorf("failed to get knowledge document for source %s %s: %w", event.Aggregate, event.AggregateID, err)
	}

	if doc.Status == "Embedded" {
		// Document might already be embedded, but since it's an update event, we should re-embed.
		log.Printf("Document %s is already Embedded, re-embedding...", doc.ID)
	}

	// 1. Delete old chunks from Milvus and PostgreSQL
	if err := p.milvusClient.DeleteVectorsByDocumentID(ctx, profile.KnowledgeCollection, doc.ID); err != nil {
		log.Printf("warning: failed to delete old vectors from milvus for doc %s: %v", doc.ID, err)
	}
	if err := p.repo.DeleteChunksByDocumentID(ctx, doc.ID); err != nil {
		return fmt.Errorf("failed to delete old chunks from db: %w", err)
	}

	// 2. Generate chunks using AI Builder with resolved provider, model, and system instruction
	aiChunks, err := p.chunkBuilder.Build(ctx, chunkProviderName, modelName, systemInstruction, doc)
	if err != nil {
		return fmt.Errorf("failed to generate chunks via ai builder: %w", err)
	}

	if len(aiChunks) == 0 {
		log.Printf("no chunks generated for document %s", doc.ID)
		return nil
	}

	var chunks []entity.KnowledgeChunk
	var vectors []milvus.KnowledgeVector

	for _, ac := range aiChunks {
		chunkID := ulid.New()
		ac.ID = chunkID
		ac.CreatedAt = time.Now()
		ac.TokenCount = 0 // can be calculated if needed

		// 3. Generate Embedding for each chunk's content using the profile's provider and model
		embeddingProvider, err := p.aiRegistry.Get(profile.Provider)
		if err != nil {
			return fmt.Errorf("failed to get embedding provider %q: %w", profile.Provider, err)
		}
		embedding, err := embeddingProvider.GenerateEmbedding(ctx, profile.Model, ac.Content)
		if err != nil {
			return fmt.Errorf("failed to generate embedding for chunk %s: %w", chunkID, err)
		}

		chunks = append(chunks, ac)

		vectors = append(vectors, milvus.KnowledgeVector{
			ChunkID:    chunkID,
			DocumentID: doc.ID,
			SourceType: doc.SourceType,
			SourceID:   doc.SourceID,
			Embedding:  embedding,
		})
	}

	// 4. Save chunks to DB
	if err := p.repo.CreateChunks(ctx, chunks); err != nil {
		return fmt.Errorf("failed to save chunks to db: %w", err)
	}

	// 5. Upsert to Milvus using active profile's collection name
	if err := p.milvusClient.UpsertVectors(ctx, profile.KnowledgeCollection, vectors); err != nil {
		return fmt.Errorf("failed to upsert vectors to milvus: %w", err)
	}

	// 6. Update document status
	now := time.Now()
	doc.Status = "Embedded"
	doc.EmbeddingModel = profile.Model
	doc.LastEmbeddedAt = &now

	if err := p.repo.UpdateDocument(ctx, doc); err != nil {
		return fmt.Errorf("failed to update document status: %w", err)
	}

	log.Printf("successfully processed and embedded document %s with %d chunks using model %s", doc.ID, len(chunks), modelName)
	return nil
}
