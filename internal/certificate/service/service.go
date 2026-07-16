// internal/certificate/service/service.go
package service

import (
	"context"
	"log"

	"dan-ai/internal/certificate/entity"
	"dan-ai/internal/certificate/repository"
	knowledgeBuilder "dan-ai/internal/knowledge/builder"
	knowledgeService "dan-ai/internal/knowledge/service"
	"dan-ai/pkg/ulid"
)

// Service defines the interface for Certificate business operations.
type Service interface {
	List(ctx context.Context, page, limit int) ([]*entity.Certificate, int64, error)
	Get(ctx context.Context, id string) (*entity.Certificate, error)
	Create(ctx context.Context, cert *entity.Certificate) error
	Update(ctx context.Context, cert *entity.Certificate) error
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

func (s *service) List(ctx context.Context, page, limit int) ([]*entity.Certificate, int64, error) {
	return s.repo.List(ctx, page, limit)
}

func (s *service) Get(ctx context.Context, id string) (*entity.Certificate, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Create(ctx context.Context, certificate *entity.Certificate) error {
	if certificate.ID == "" {
		certificate.ID = ulid.New()
	}
	err := s.repo.Create(ctx, certificate)
	if err == nil && s.knowledgeSvc != nil {
		go func(c entity.Certificate) {
			title, content := knowledgeBuilder.BuildCertificateDocument(c)
			if err := s.knowledgeSvc.Sync(context.Background(), "certificate", c.ID, title, content); err != nil {
				log.Printf("Failed to sync certificate knowledge: %v\n", err)
			}
		}(*certificate)
	}
	return err
}

func (s *service) Update(ctx context.Context, certificate *entity.Certificate) error {
	err := s.repo.Update(ctx, certificate)
	if err == nil && s.knowledgeSvc != nil {
		go func(c entity.Certificate) {
			title, content := knowledgeBuilder.BuildCertificateDocument(c)
			if err := s.knowledgeSvc.Sync(context.Background(), "certificate", c.ID, title, content); err != nil {
				log.Printf("Failed to sync certificate knowledge: %v\n", err)
			}
		}(*certificate)
	}
	return err
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
