package chunk

import (
	"context"
	"fmt"
	"portfolio-ai/internal/ai/provider"
	"portfolio-ai/internal/knowledge/entity"
)

type Builder interface {
	Build(ctx context.Context, document *entity.KnowledgeDocument) ([]entity.KnowledgeChunk, error)
}

type aiBuilder struct {
	aiProvider provider.Provider
}

func NewAIBuilder(aiProvider provider.Provider) Builder {
	return &aiBuilder{
		aiProvider: aiProvider,
	}
}

func (b *aiBuilder) Build(ctx context.Context, document *entity.KnowledgeDocument) ([]entity.KnowledgeChunk, error) {
	// Generate chunks from LLM
	aiChunks, err := b.aiProvider.GenerateChunks(ctx, document.Title, document.Content)
	if err != nil {
		return nil, fmt.Errorf("ai provider failed to generate chunks: %w", err)
	}

	var chunks []entity.KnowledgeChunk
	for i, ac := range aiChunks {
		// Format the chunk into the structure our DB accepts
		embedText := fmt.Sprintf("Title: %s\nContent: %s\nKeywords: %v", ac.Title, ac.Content, ac.Keywords)
		
		chunks = append(chunks, entity.KnowledgeChunk{
			// ID will be assigned by the caller (processor)
			DocumentID:     document.ID,
			ChunkIndex:     int32(i),
			Content:        embedText,
			EmbeddingModel: document.EmbeddingModel,
		})
	}
	return chunks, nil
}
