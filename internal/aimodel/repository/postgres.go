// internal/aimodel/repository/postgres.go
package repository

import (
	"context"
	"portfolio-ai/internal/aimodel/entity"

	"gorm.io/gorm"
)

// Repository defines the interface for AIModel database operations.
type Repository interface {
	Get(ctx context.Context, id string) (*entity.AIModel, error)
	List(ctx context.Context, enabledOnly bool) ([]entity.AIModel, error)
	Create(ctx context.Context, model *entity.AIModel) error
	Update(ctx context.Context, model *entity.AIModel) error
	Delete(ctx context.Context, id string) error
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new Repository implementation using GORM.
func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Get(ctx context.Context, id string) (*entity.AIModel, error) {
	var model entity.AIModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *postgresRepository) List(ctx context.Context, enabledOnly bool) ([]entity.AIModel, error) {
	var models []entity.AIModel
	q := r.db.WithContext(ctx).Order("name ASC")
	if enabledOnly {
		q = q.Where("enabled = ?", true)
	}
	if err := q.Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *postgresRepository) Create(ctx context.Context, model *entity.AIModel) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *postgresRepository) Update(ctx context.Context, model *entity.AIModel) error {
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.AIModel{}, "id = ?", id).Error
}
