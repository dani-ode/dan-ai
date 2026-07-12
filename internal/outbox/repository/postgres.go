package repository

import (
	"context"
	"time"

	"portfolio-ai/internal/outbox/entity"

	"gorm.io/gorm"
)

type Repository interface {
	GetUnpublishedEvents(ctx context.Context, limit int) ([]entity.OutboxEvent, error)
	MarkAsPublished(ctx context.Context, id string) error
	MarkAsFailed(ctx context.Context, id string, reason string) error
}

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetUnpublishedEvents(ctx context.Context, limit int) ([]entity.OutboxEvent, error) {
	var events []entity.OutboxEvent
	// We get events where published is false and retry_count < 3
	err := r.db.WithContext(ctx).
		Where("published = ? AND retry_count < ?", false, 3).
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

func (r *postgresRepository) MarkAsPublished(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"published":    true,
			"published_at": &now,
		}).Error
}

func (r *postgresRepository) MarkAsFailed(ctx context.Context, id string, reason string) error {
	return r.db.WithContext(ctx).
		Model(&entity.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"retry_count":   gorm.Expr("retry_count + 1"),
			"failed_reason": reason,
		}).Error
}
