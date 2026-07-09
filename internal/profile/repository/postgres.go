// internal/profile/repository/postgres.go
package repository

import (
	"context"
	"portfolio-ai/internal/profile/entity"

	"gorm.io/gorm"
)

// Repository defines the interface for Profile database operations.
type Repository interface {
	Get(ctx context.Context, id string) (*entity.Profile, error)
	Create(ctx context.Context, profile *entity.Profile) error
	Update(ctx context.Context, profile *entity.Profile) error
	Delete(ctx context.Context, id string) error
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new Repository implementation using GORM.
func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Get(ctx context.Context, id string) (*entity.Profile, error) {
	var profile entity.Profile
	if err := r.db.WithContext(ctx).First(&profile, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *postgresRepository) Create(ctx context.Context, profile *entity.Profile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

func (r *postgresRepository) Update(ctx context.Context, profile *entity.Profile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Profile{}, "id = ?", id).Error
}
