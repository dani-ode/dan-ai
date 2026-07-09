// internal/prompt/service/service.go
package service

import (
	"context"
	"portfolio-ai/internal/prompt/entity"
	"portfolio-ai/internal/prompt/repository"
	"portfolio-ai/pkg/ulid"
)

// Service defines the interface for Prompt business operations.
type Service interface {
	Get(ctx context.Context, id string) (*entity.Prompt, error)
	List(ctx context.Context, activeOnly bool) ([]entity.Prompt, error)
	Create(ctx context.Context, prompt *entity.Prompt) error
	Update(ctx context.Context, prompt *entity.Prompt) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo repository.Repository
}

// NewService creates a new Service instance.
func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get(ctx context.Context, id string) (*entity.Prompt, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) List(ctx context.Context, activeOnly bool) ([]entity.Prompt, error) {
	return s.repo.List(ctx, activeOnly)
}

func (s *service) Create(ctx context.Context, prompt *entity.Prompt) error {
	if prompt.ID == "" {
		prompt.ID = ulid.New()
	}
	if prompt.Version == 0 {
		prompt.Version = 1
	}
	return s.repo.Create(ctx, prompt)
}

func (s *service) Update(ctx context.Context, prompt *entity.Prompt) error {
	return s.repo.Update(ctx, prompt)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
