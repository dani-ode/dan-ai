package chunk

import (
	"context"
	"fmt"
	"dan-ai/internal/ai/provider"
	"dan-ai/internal/knowledge/entity"
)

type Builder interface {
	Build(ctx context.Context, providerName, modelName, systemInstruction string, document *entity.KnowledgeDocument) ([]entity.KnowledgeChunk, error)
}

type aiBuilder struct {
	aiRegistry *provider.Registry
}

func NewAIBuilder(aiRegistry *provider.Registry) Builder {
	return &aiBuilder{
		aiRegistry: aiRegistry,
	}
}

func (b *aiBuilder) Build(ctx context.Context, providerName, modelName, systemInstruction string, document *entity.KnowledgeDocument) ([]entity.KnowledgeChunk, error) {
	// Resolve provider dynamically
	aiProvider, err := b.aiRegistry.Get(providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider %q for chunk building: %w", providerName, err)
	}

	// Generate chunks from LLM
	aiChunks, err := aiProvider.GenerateChunks(ctx, modelName, systemInstruction, document.Title, document.Content)
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

