// internal/experience/service/service.go
package service

import (
	"context"
	"log"

	"dan-ai/internal/experience/entity"
	"dan-ai/internal/experience/repository"
	knowledgeBuilder "dan-ai/internal/knowledge/builder"
	knowledgeService "dan-ai/internal/knowledge/service"
	"dan-ai/pkg/ulid"
)

// Service defines the interface for Experience business operations.
type Service interface {
	List(ctx context.Context, page, limit int) ([]*entity.Experience, int64, error)
	Get(ctx context.Context, id string) (*entity.Experience, error)
	Create(ctx context.Context, experience *entity.Experience) error
	Update(ctx context.Context, experience *entity.Experience) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo         repository.Repository
	knowledgeSvc knowledgeService.Service
}

// NewService creates a new Service instance.
func NewService(repo repository.Repository, knowledgeSvc knowledgeService.Service) Service {
	return &service{
		repo:         repo,
		knowledgeSvc: knowledgeSvc,
	}
}

func (s *service) List(ctx context.Context, page, limit int) ([]*entity.Experience, int64, error) {
	return s.repo.List(ctx, page, limit)
}

func (s *service) Get(ctx context.Context, id string) (*entity.Experience, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Create(ctx context.Context, experience *entity.Experience) error {
	if experience.ID == "" {
		experience.ID = ulid.New()
	}
	err := s.repo.Create(ctx, experience)
	if err == nil && s.knowledgeSvc != nil {
		go func(e entity.Experience) {
			title, content := knowledgeBuilder.BuildExperienceDocument(e)
			if err := s.knowledgeSvc.Sync(context.Background(), "experience", e.ID, title, content); err != nil {
				log.Printf("Failed to sync experience knowledge: %v\n", err)
			}
		}(*experience)
	}
	return err
}

func (s *service) Update(ctx context.Context, experience *entity.Experience) error {
	err := s.repo.Update(ctx, experience)
	if err == nil && s.knowledgeSvc != nil {
		go func(e entity.Experience) {
			title, content := knowledgeBuilder.BuildExperienceDocument(e)
			if err := s.knowledgeSvc.Sync(context.Background(), "experience", e.ID, title, content); err != nil {
				log.Printf("Failed to sync experience knowledge: %v\n", err)
			}
		}(*experience)
	}
	return err
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
