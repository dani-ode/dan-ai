package service

import (
	"context"

	"dan-ai/internal/memory/entity"
	"dan-ai/internal/memory/repository"
)

type Service interface {
	SaveMemories(ctx context.Context, memories []entity.Memory) error
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) SaveMemories(ctx context.Context, memories []entity.Memory) error {
	for _, memory := range memories {
		if err := s.repo.UpsertMemory(ctx, &memory); err != nil {
			return err
		}
	}
	return nil
}
