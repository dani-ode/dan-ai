// internal/profile/service/service.go
package service

import (
	"context"
	"log"

	knowledgeBuilder "dan-ai/internal/knowledge/builder"
	knowledgeService "dan-ai/internal/knowledge/service"
	"dan-ai/internal/profile/entity"
	"dan-ai/internal/profile/repository"
	"dan-ai/pkg/ulid"
)

// Service defines the interface for Profile business operations.
type Service interface {
	Get(ctx context.Context, id string) (*entity.Profile, error)
	Create(ctx context.Context, profile *entity.Profile) error
	Update(ctx context.Context, profile *entity.Profile) error
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

func (s *service) Get(ctx context.Context, id string) (*entity.Profile, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Create(ctx context.Context, profile *entity.Profile) error {
	if profile.ID == "" {
		profile.ID = ulid.New()
	}
	err := s.repo.Create(ctx, profile)
	if err == nil && s.knowledgeSvc != nil {
		go func(p entity.Profile) {
			title, content := knowledgeBuilder.BuildProfileDocument(p)
			if err := s.knowledgeSvc.Sync(context.Background(), "profile", p.ID, title, content); err != nil {
				log.Printf("Failed to sync profile knowledge: %v\n", err)
			}
		}(*profile)
	}
	return err
}

func (s *service) Update(ctx context.Context, profile *entity.Profile) error {
	err := s.repo.Update(ctx, profile)
	if err == nil && s.knowledgeSvc != nil {
		go func(p entity.Profile) {
			title, content := knowledgeBuilder.BuildProfileDocument(p)
			if err := s.knowledgeSvc.Sync(context.Background(), "profile", p.ID, title, content); err != nil {
				log.Printf("Failed to sync profile knowledge: %v\n", err)
			}
		}(*profile)
	}
	return err
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
