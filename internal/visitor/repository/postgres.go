// internal/visitor/repository/postgres.go
package repository

import (
	"context"
	"dan-ai/internal/visitor/entity"
	"time"

	"gorm.io/gorm"
)

// Repository defines the interface for Visitor database operations.
type Repository interface {
	Get(ctx context.Context, id string) (*entity.Visitor, error)
	Create(ctx context.Context, visitor *entity.Visitor) error
	UpdateLastSeen(ctx context.Context, id string) error
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new Repository implementation using GORM.
func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Get(ctx context.Context, id string) (*entity.Visitor, error) {
	var visitor entity.Visitor
	if err := r.db.WithContext(ctx).First(&visitor, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &visitor, nil
}

func (r *postgresRepository) Create(ctx context.Context, visitor *entity.Visitor) error {
	return r.db.WithContext(ctx).Create(visitor).Error
}

func (r *postgresRepository) UpdateLastSeen(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&entity.Visitor{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"last_seen_at": time.Now(),
			"updated_at":   time.Now(),
		}).Error
}
