// internal/prompt/repository/postgres.go
package repository

import (
	"context"
	"dan-ai/internal/prompt/entity"

	"gorm.io/gorm"
)

// Repository defines the interface for Prompt database operations.
type Repository interface {
	Get(ctx context.Context, id string) (*entity.Prompt, error)
	List(ctx context.Context, activeOnly bool) ([]entity.Prompt, error)
	Create(ctx context.Context, prompt *entity.Prompt) error
	Update(ctx context.Context, prompt *entity.Prompt) error
	Delete(ctx context.Context, id string) error
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new Repository implementation using GORM.
func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Get(ctx context.Context, id string) (*entity.Prompt, error) {
	var prompt entity.Prompt
	if err := r.db.WithContext(ctx).Preload("AIModel").First(&prompt, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &prompt, nil
}


func (r *postgresRepository) List(ctx context.Context, activeOnly bool) ([]entity.Prompt, error) {
	var prompts []entity.Prompt
	q := r.db.WithContext(ctx).Order("created_at DESC")
	if activeOnly {
		q = q.Where("active = ?", true)
	}
	if err := q.Find(&prompts).Error; err != nil {
		return nil, err
	}
	return prompts, nil
}

func (r *postgresRepository) Create(ctx context.Context, prompt *entity.Prompt) error {
	return r.db.WithContext(ctx).Create(prompt).Error
}

func (r *postgresRepository) Update(ctx context.Context, prompt *entity.Prompt) error {
	return r.db.WithContext(ctx).Save(prompt).Error
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Prompt{}, "id = ?", id).Error
}
