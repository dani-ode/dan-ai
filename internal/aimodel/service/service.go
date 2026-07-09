// internal/aimodel/service/service.go
package service

import (
	"context"
	"portfolio-ai/internal/aimodel/entity"
	"portfolio-ai/internal/aimodel/repository"
	"portfolio-ai/pkg/ulid"
)

// Service defines the interface for AIModel business operations.
type Service interface {
	Get(ctx context.Context, id string) (*entity.AIModel, error)
	List(ctx context.Context, enabledOnly bool) ([]entity.AIModel, error)
	Create(ctx context.Context, model *entity.AIModel) error
	Update(ctx context.Context, model *entity.AIModel) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo repository.Repository
}

// NewService creates a new Service instance.
func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get(ctx context.Context, id string) (*entity.AIModel, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) List(ctx context.Context, enabledOnly bool) ([]entity.AIModel, error) {
	return s.repo.List(ctx, enabledOnly)
}

func (s *service) Create(ctx context.Context, model *entity.AIModel) error {
	if model.ID == "" {
		model.ID = ulid.New()
	}
	return s.repo.Create(ctx, model)
}

func (s *service) Update(ctx context.Context, model *entity.AIModel) error {
	return s.repo.Update(ctx, model)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
