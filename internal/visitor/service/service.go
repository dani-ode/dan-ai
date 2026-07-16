// internal/visitor/service/service.go
package service

import (
	"context"
	"errors"
	"dan-ai/internal/visitor/entity"
	"dan-ai/internal/visitor/repository"
	"dan-ai/pkg/ulid"

	"gorm.io/gorm"
)

// Service defines the interface for Visitor business operations.
type Service interface {
	Register(ctx context.Context, visitorID string) (*entity.Visitor, error)
	Get(ctx context.Context, id string) (*entity.Visitor, error)
}

type service struct {
	repo repository.Repository
}

// NewService creates a new Service instance.
func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

// Register performs an upsert: if visitor_id is provided and exists, update last_seen_at.
// If visitor_id is empty or doesn't exist, create a new visitor.
func (s *service) Register(ctx context.Context, visitorID string) (*entity.Visitor, error) {
	// If visitor_id is provided, try to find existing visitor
	if visitorID != "" {
		existing, err := s.repo.Get(ctx, visitorID)
		if err == nil {
			// Visitor exists, update last_seen_at
			if updateErr := s.repo.UpdateLastSeen(ctx, visitorID); updateErr != nil {
				return nil, updateErr
			}
			// Re-fetch to get updated timestamps
			return s.repo.Get(ctx, visitorID)
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		// Not found — fall through to create with the provided ID
		_ = existing
	}

	// Create new visitor
	visitor := &entity.Visitor{
		ID: visitorID,
	}
	if visitor.ID == "" {
		visitor.ID = ulid.New()
	}

	if err := s.repo.Create(ctx, visitor); err != nil {
		return nil, err
	}
	return visitor, nil
}

func (s *service) Get(ctx context.Context, id string) (*entity.Visitor, error) {
	return s.repo.Get(ctx, id)
}
