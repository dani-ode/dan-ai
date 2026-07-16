package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"dan-ai/internal/knowledge/entity"
	"dan-ai/internal/knowledge/repository"
	outboxEntity "dan-ai/internal/outbox/entity"

	"github.com/oklog/ulid/v2"
)

type Service interface {
	Sync(ctx context.Context, sourceType, sourceID, title, content string) error
	// For gRPC Handlers:
	GetDocument(ctx context.Context, id string) (*entity.KnowledgeDocument, error)
	ListDocuments(ctx context.Context, page, pageSize int, sourceType string) ([]entity.KnowledgeDocument, int64, error)
	ListChunks(ctx context.Context, documentID string) ([]entity.KnowledgeChunk, error)
}

type service struct {
	repo repository.KnowledgeRepository
}

func NewService(repo repository.KnowledgeRepository) Service {
	return &service{repo: repo}
}

func (s *service) Sync(ctx context.Context, sourceType, sourceID, title, content string) error {
	// 1. Generate Checksum
	hash := sha256.Sum256([]byte(content))
	checksum := hex.EncodeToString(hash[:])

	// 2. Check if Document exists
	existingDoc, err := s.repo.GetDocumentBySource(ctx, sourceType, sourceID)
	if err != nil && err.Error() != "record not found" {
		return fmt.Errorf("failed to get existing document: %w", err)
	}

	docID := ""
	version := int32(1)

	if existingDoc != nil {
		if existingDoc.Checksum == checksum {
			// No changes needed
			return nil
		}
		// Update existing
		existingDoc.Title = title
		existingDoc.Content = content
		existingDoc.Checksum = checksum
		existingDoc.Version += 1
		existingDoc.Status = "Pending"
		existingDoc.UpdatedAt = time.Now()
		
		if err := s.repo.UpdateDocument(ctx, existingDoc); err != nil {
			return fmt.Errorf("failed to update document: %w", err)
		}
		docID = existingDoc.ID
		version = existingDoc.Version

		// Delete old chunks
		if err := s.repo.DeleteChunksByDocumentID(ctx, docID); err != nil {
			return fmt.Errorf("failed to delete old chunks: %w", err)
		}
	} else {
		// Create new document
		docID = ulid.Make().String()
		newDoc := &entity.KnowledgeDocument{
			ID:         docID,
			SourceType: sourceType,
			SourceID:   sourceID,
			Title:      title,
			Content:    content,
			Checksum:   checksum,
			Version:    version,
			Status:     "Pending",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if err := s.repo.CreateDocument(ctx, newDoc); err != nil {
			return fmt.Errorf("failed to create document: %w", err)
		}
	}

	// 3. Create OutboxEvent to trigger the AI Embedding Worker
	outboxEvent := &outboxEntity.OutboxEvent{
		ID:          ulid.Make().String(),
		Aggregate:   sourceType,
		AggregateID: sourceID,
		EventType:   "updated",
		Payload:     []byte(`{"source_type":"` + sourceType + `","source_id":"` + sourceID + `"}`),
		Published:   false,
		RetryCount:  0,
		CreatedAt:   time.Now(),
	}
	if err := s.repo.CreateOutboxEvent(ctx, outboxEvent); err != nil {
		return fmt.Errorf("failed to create outbox event: %w", err)
	}

	return nil
}

func (s *service) GetDocument(ctx context.Context, id string) (*entity.KnowledgeDocument, error) {
	return s.repo.GetDocumentByID(ctx, id)
}

func (s *service) ListDocuments(ctx context.Context, page, pageSize int, sourceType string) ([]entity.KnowledgeDocument, int64, error) {
	return s.repo.ListDocuments(ctx, page, pageSize, sourceType)
}

func (s *service) ListChunks(ctx context.Context, documentID string) ([]entity.KnowledgeChunk, error) {
	return s.repo.ListChunksByDocumentID(ctx, documentID)
}
