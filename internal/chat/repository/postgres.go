// internal/chat/repository/postgres.go
package repository

import (
	"context"
	"dan-ai/internal/chat/entity"
	"time"

	"gorm.io/gorm"
)

// Repository defines the interface for Chat database operations (sessions and messages).
type Repository interface {
	// Session operations
	CreateSession(ctx context.Context, session *entity.ChatSession) error
	GetSession(ctx context.Context, id string) (*entity.ChatSession, error)
	ListSessionsByVisitor(ctx context.Context, visitorID string) ([]entity.ChatSession, error)
	RenameSession(ctx context.Context, id, title string) (*entity.ChatSession, error)
	DeleteSession(ctx context.Context, id string) error

	// Message operations
	CreateMessage(ctx context.Context, message *entity.ChatMessage) error
	GetMessageByID(ctx context.Context, id string) (*entity.ChatMessage, error)
	ListMessagesBySession(ctx context.Context, sessionID string) ([]entity.ChatMessage, error)
	DeleteMessage(ctx context.Context, id string) error
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new Repository implementation using GORM.
func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

// --- Session operations ---

func (r *postgresRepository) CreateSession(ctx context.Context, session *entity.ChatSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *postgresRepository) GetSession(ctx context.Context, id string) (*entity.ChatSession, error) {
	var session entity.ChatSession
	if err := r.db.WithContext(ctx).First(&session, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *postgresRepository) ListSessionsByVisitor(ctx context.Context, visitorID string) ([]entity.ChatSession, error) {
	var sessions []entity.ChatSession
	if err := r.db.WithContext(ctx).
		Where("visitor_id = ?", visitorID).
		Order("created_at DESC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *postgresRepository) RenameSession(ctx context.Context, id, title string) (*entity.ChatSession, error) {
	result := r.db.WithContext(ctx).
		Model(&entity.ChatSession{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"title":      title,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return r.GetSession(ctx, id)
}

func (r *postgresRepository) DeleteSession(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&entity.ChatSession{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// --- Message operations ---

func (r *postgresRepository) CreateMessage(ctx context.Context, message *entity.ChatMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *postgresRepository) ListMessagesBySession(ctx context.Context, sessionID string) ([]entity.ChatMessage, error) {
	var messages []entity.ChatMessage
	if err := r.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *postgresRepository) GetMessageByID(ctx context.Context, id string) (*entity.ChatMessage, error) {
	var message entity.ChatMessage
	if err := r.db.WithContext(ctx).First(&message, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *postgresRepository) DeleteMessage(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&entity.ChatMessage{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
