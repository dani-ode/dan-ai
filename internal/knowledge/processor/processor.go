package processor

import (
	"context"
	"fmt"
	"log"
	"portfolio-ai/internal/ai/provider"
	"portfolio-ai/internal/knowledge/chunk"
	"portfolio-ai/internal/knowledge/entity"
	"portfolio-ai/internal/knowledge/repository"
	"portfolio-ai/pkg/kafka"
	"portfolio-ai/pkg/milvus"
	"portfolio-ai/pkg/ulid"
	"time"
)

type Processor struct {
	repo         repository.KnowledgeRepository
	aiProvider   provider.Provider
	milvusClient *milvus.Client
	chunkBuilder chunk.Builder
}

func NewProcessor(repo repository.KnowledgeRepository, aiProvider provider.Provider, milvusClient *milvus.Client, chunkBuilder chunk.Builder) *Processor {
	return &Processor{
		repo:         repo,
		aiProvider:   aiProvider,
		milvusClient: milvusClient,
		chunkBuilder: chunkBuilder,
	}
}

func (p *Processor) ProcessEvent(ctx context.Context, event kafka.Event) error {
	log.Printf("processing knowledge event for %s %s", event.Aggregate, event.AggregateID)

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
	if err := p.milvusClient.DeleteVectorsByDocumentID(ctx, doc.ID); err != nil {
		log.Printf("warning: failed to delete old vectors from milvus for doc %s: %v", doc.ID, err)
	}
	if err := p.repo.DeleteChunksByDocumentID(ctx, doc.ID); err != nil {
		return fmt.Errorf("failed to delete old chunks from db: %w", err)
	}

	// 2. Generate chunks using AI Builder
	aiChunks, err := p.chunkBuilder.Build(ctx, doc)
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

		// 3. Generate Embedding for each chunk's content
		embedding, err := p.aiProvider.GenerateEmbedding(ctx, ac.Content)
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

	// 5. Upsert to Milvus
	if err := p.milvusClient.UpsertVectors(ctx, vectors); err != nil {
		return fmt.Errorf("failed to upsert vectors to milvus: %w", err)
	}

	// 6. Update document status
	now := time.Now()
	doc.Status = "Embedded"
	doc.EmbeddingModel = "gemini-embedding-2"
	doc.LastEmbeddedAt = &now

	if err := p.repo.UpdateDocument(ctx, doc); err != nil {
		return fmt.Errorf("failed to update document status: %w", err)
	}

	log.Printf("successfully processed and embedded document %s with %d chunks", doc.ID, len(chunks))
	return nil
}
