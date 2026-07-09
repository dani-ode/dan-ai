// internal/profile/service/service.go
package service

import (
	"context"
	"portfolio-ai/internal/profile/entity"
	"portfolio-ai/internal/profile/repository"
	"portfolio-ai/pkg/ulid"
)

// Service defines the interface for Profile business operations.
type Service interface {
	Get(ctx context.Context, id string) (*entity.Profile, error)
	Create(ctx context.Context, profile *entity.Profile) error
	Update(ctx context.Context, profile *entity.Profile) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo repository.Repository
}

// NewService creates a new Service instance.
func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get(ctx context.Context, id string) (*entity.Profile, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Create(ctx context.Context, profile *entity.Profile) error {
	if profile.ID == "" {
		profile.ID = ulid.New()
	}
	return s.repo.Create(ctx, profile)
}

func (s *service) Update(ctx context.Context, profile *entity.Profile) error {
	return s.repo.Update(ctx, profile)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
