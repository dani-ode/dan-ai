// internal/skill/service/service.go
package service

import (
	"context"
	"log"

	knowledgeBuilder "portfolio-ai/internal/knowledge/builder"
	knowledgeService "portfolio-ai/internal/knowledge/service"
	"portfolio-ai/internal/skill/entity"
	"portfolio-ai/internal/skill/repository"
	"portfolio-ai/pkg/ulid"
)

// Service defines the interface for Skill business operations.
type Service interface {
	List(ctx context.Context, page, limit int) ([]*entity.Skill, int64, error)
	Get(ctx context.Context, id string) (*entity.Skill, error)
	Create(ctx context.Context, skill *entity.Skill) error
	Update(ctx context.Context, skill *entity.Skill) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo         repository.Repository
	knowledgeSvc knowledgeService.Service
}

func NewService(repo repository.Repository, knowledgeSvc knowledgeService.Service) Service {
	return &service{
		repo:         repo,
		knowledgeSvc: knowledgeSvc,
	}
}

func (s *service) List(ctx context.Context, page, limit int) ([]*entity.Skill, int64, error) {
	return s.repo.List(ctx, page, limit)
}

func (s *service) Get(ctx context.Context, id string) (*entity.Skill, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Create(ctx context.Context, skill *entity.Skill) error {
	if skill.ID == "" {
		skill.ID = ulid.New()
	}
	err := s.repo.Create(ctx, skill)
	if err == nil && s.knowledgeSvc != nil {
		go func(sk entity.Skill) {
			title, content := knowledgeBuilder.BuildSkillDocument(sk)
			if err := s.knowledgeSvc.Sync(context.Background(), "skill", sk.ID, title, content); err != nil {
				log.Printf("Failed to sync skill knowledge: %v\n", err)
			}
		}(*skill)
	}
	return err
}

func (s *service) Update(ctx context.Context, skill *entity.Skill) error {
	err := s.repo.Update(ctx, skill)
	if err == nil && s.knowledgeSvc != nil {
		go func(sk entity.Skill) {
			title, content := knowledgeBuilder.BuildSkillDocument(sk)
			if err := s.knowledgeSvc.Sync(context.Background(), "skill", sk.ID, title, content); err != nil {
				log.Printf("Failed to sync skill knowledge: %v\n", err)
			}
		}(*skill)
	}
	return err
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
