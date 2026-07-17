package repository

import (
	"context"
	"time"

	"dan-ai/internal/memory/entity"

	"gorm.io/gorm"
)

type Repository interface {
	UpsertMemory(ctx context.Context, memory *entity.Memory) error
	ListByVisitor(ctx context.Context, visitorID string) ([]entity.Memory, error)
	GetMemoriesByIDs(ctx context.Context, ids []string) ([]entity.Memory, error)
}

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) UpsertMemory(ctx context.Context, memory *entity.Memory) error {
	now := time.Now()
	memory.UpdatedAt = now
	if memory.CreatedAt.IsZero() {
		memory.CreatedAt = now
	}

	return r.db.WithContext(ctx).Save(memory).Error
}

func (r *postgresRepository) ListByVisitor(ctx context.Context, visitorID string) ([]entity.Memory, error) {
	var memories []entity.Memory
	err := r.db.WithContext(ctx).
		Where("visitor_id = ?", visitorID).
		Order("updated_at DESC").
		Find(&memories).Error
	return memories, err
}

func (r *postgresRepository) GetMemoriesByIDs(ctx context.Context, ids []string) ([]entity.Memory, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var memories []entity.Memory
	err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&memories).Error
	return memories, err
}
