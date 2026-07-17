package service

import (
	"context"
	"fmt"

	embeddingrepo "dan-ai/internal/embedding/repository"
	"dan-ai/internal/memory/entity"
	"dan-ai/internal/memory/repository"
	aiprovider "dan-ai/internal/ai/provider"
	promptrepo "dan-ai/internal/prompt/repository"
	"dan-ai/pkg/milvus"
)

type Service interface {
	SaveMemories(ctx context.Context, modelName string, memories []entity.Memory) error
}

type service struct {
	repo          repository.Repository
	milvusClient  *milvus.Client
	aiRegistry    *aiprovider.Registry
	promptRepo    promptrepo.Repository
	embeddingRepo embeddingrepo.Repository
}

func NewService(repo repository.Repository, milvusClient *milvus.Client, aiRegistry *aiprovider.Registry, promptRepo promptrepo.Repository, embeddingRepo embeddingrepo.Repository) Service {
	return &service{
		repo:          repo,
		milvusClient:  milvusClient,
		aiRegistry:    aiRegistry,
		promptRepo:    promptRepo,
		embeddingRepo: embeddingRepo,
	}
}

const MergeSystemInstruction = "You are a memory consolidator. Your task is to combine two similar visitor memories into a single, cohesive memory without losing key context. Keep it short and concise (max 250 characters). Do not add metadata or introductions. Return only the consolidated memory string."

func (s *service) SaveMemories(ctx context.Context, modelName string, memories []entity.Memory) error {
	// Get active embedding profile
	profile, err := s.embeddingRepo.GetActiveProfile(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active embedding profile: %w", err)
	}

	// Resolve dynamic system prompt from database
	systemInstruction := MergeSystemInstruction
	allPrompts, err := s.promptRepo.List(ctx, false)
	if err == nil {
		for _, p := range allPrompts {
			if p.Name == "Memory Consolidator" {
				systemInstruction = p.SystemPrompt
				break
			}
		}
	}

	for _, memory := range memories {
		// 1. Generate embedding for the new memory using profile's provider and model
		embeddingProvider, err := s.aiRegistry.Get(profile.Provider)
		if err != nil {
			return fmt.Errorf("failed to get embedding provider %q: %w", profile.Provider, err)
		}
		embedding, err := embeddingProvider.GenerateEmbedding(ctx, profile.Model, memory.MemoryText)
		if err != nil {
			return fmt.Errorf("failed to generate embedding: %w", err)
		}

		// 2. Search Milvus for similar existing memories of the visitor (top 3)
		similarMemories, err := s.milvusClient.SearchVisitorMemory(ctx, profile.VisitorCollection, memory.VisitorID, embedding, 3)
		if err != nil {
			return fmt.Errorf("failed to search similar memories: %w", err)
		}

		var targetMemory entity.Memory = memory
		var finalEmbedding []float32 = embedding

		// 3. Consolidate if top match is similar enough (score >= 0.85)
		if len(similarMemories) > 0 && similarMemories[0].Score >= 0.85 {
			bestMatch := similarMemories[0]

			// Fetch the existing memory text from PostgreSQL (since Milvus no longer stores it)
			existingMemList, err := s.repo.GetMemoriesByIDs(ctx, []string{bestMatch.MemoryID})
			if err != nil || len(existingMemList) == 0 {
				return fmt.Errorf("failed to fetch existing memory text: %w", err)
			}
			existingMemory := existingMemList[0]

			// Call LLM to combine/merge using the model's provider
			prompt := fmt.Sprintf("Memory A:\n%s\n\nMemory B:\n%s\n\nCombine into one consolidated memory:", existingMemory.MemoryText, memory.MemoryText)
			chatProvider, err := s.aiRegistry.Get("gemini") // default to gemini for consolidation
			if err != nil {
				return fmt.Errorf("failed to get chat provider for consolidation: %w", err)
			}
			mergedResp, err := chatProvider.GenerateChatResponse(ctx, modelName, systemInstruction, prompt)
			if err != nil {
				return fmt.Errorf("failed to merge memories: %w", err)
			}

			// Clean up merged response
			mergedText := mergedResp.Content
			// Re-embed the merged text
			mergedEmbedding, err := embeddingProvider.GenerateEmbedding(ctx, profile.Model, mergedText)
			if err != nil {
				return fmt.Errorf("failed to generate merged embedding: %w", err)
			}

			targetMemory.ID = bestMatch.MemoryID
			targetMemory.MemoryText = mergedText
			finalEmbedding = mergedEmbedding
		}

		// 4. Save to PostgreSQL
		if err := s.repo.UpsertMemory(ctx, &targetMemory); err != nil {
			return fmt.Errorf("failed to save memory to postgres: %w", err)
		}

		// 5. Upsert to Milvus using active profile's collection name
		vector := milvus.VisitorMemoryVector{
			MemoryID:  targetMemory.ID,
			VisitorID: targetMemory.VisitorID,
			Embedding: finalEmbedding,
		}
		if err := s.milvusClient.UpsertVisitorMemoryVectors(ctx, profile.VisitorCollection, []milvus.VisitorMemoryVector{vector}); err != nil {
			return fmt.Errorf("failed to save memory vector to milvus: %w", err)
		}
	}
	return nil
}
