// internal/chat/service/service.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"dan-ai/internal/chat/entity"
	"dan-ai/internal/chat/repository"
	outboxEntity "dan-ai/internal/outbox/entity"
	outboxrepo "dan-ai/internal/outbox/repository"
	"dan-ai/pkg/ulid"
)

// Service defines the interface for Chat business operations.
type Service interface {
	// Session operations
	CreateSession(ctx context.Context, visitorID, promptID string) (*entity.ChatSession, error)
	GetSession(ctx context.Context, id string) (*entity.ChatSession, error)
	ListSessions(ctx context.Context, visitorID string) ([]entity.ChatSession, error)
	RenameSession(ctx context.Context, id, title string) (*entity.ChatSession, error)
	DeleteSession(ctx context.Context, id string) error

	// Message operations
	CreateMessage(ctx context.Context, sessionID, role, content string) (*entity.ChatMessage, error)
	ListMessages(ctx context.Context, sessionID string) ([]entity.ChatMessage, error)
	DeleteMessage(ctx context.Context, id string) error
}

type service struct {
	repo       repository.Repository
	outboxRepo outboxrepo.Repository
}

// NewService creates a new Service instance.
func NewService(repo repository.Repository, outboxRepo outboxrepo.Repository) Service {
	return &service{repo: repo, outboxRepo: outboxRepo}
}

// --- Session operations ---

func (s *service) CreateSession(ctx context.Context, visitorID, promptID string) (*entity.ChatSession, error) {
	session := &entity.ChatSession{
		ID:        ulid.New(),
		VisitorID: visitorID,
		PromptID:  promptID,
		Title:     "New Chat",
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *service) GetSession(ctx context.Context, id string) (*entity.ChatSession, error) {
	return s.repo.GetSession(ctx, id)
}

func (s *service) ListSessions(ctx context.Context, visitorID string) ([]entity.ChatSession, error) {
	return s.repo.ListSessionsByVisitor(ctx, visitorID)
}

func (s *service) RenameSession(ctx context.Context, id, title string) (*entity.ChatSession, error) {
	return s.repo.RenameSession(ctx, id, title)
}

func (s *service) DeleteSession(ctx context.Context, id string) error {
	return s.repo.DeleteSession(ctx, id)
}

// --- Message operations ---

func (s *service) CreateMessage(ctx context.Context, sessionID, role, content string) (*entity.ChatMessage, error) {
	message := &entity.ChatMessage{
		ID:        ulid.New(),
		SessionID: sessionID,
		Role:      role,
		Content:   content,
		Status:    "Pending",
	}

	if err := s.repo.CreateMessage(ctx, message); err != nil {
		return nil, err
	}

	if role == "assistant" {
		session, err := s.repo.GetSession(ctx, sessionID)
		if err != nil {
			return nil, err
		}

		payload, err := json.Marshal(map[string]string{
			"visitor_id":          session.VisitorID,
			"session_id":          sessionID,
			"assistant_message_id": message.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal chat.completed payload: %w", err)
		}

		outboxEvent := &outboxEntity.OutboxEvent{
			ID:          ulid.New(),
			Aggregate:   "chat_session",
			AggregateID: sessionID,
			EventType:   "chat.completed",
			Payload:     payload,
			Published:   false,
			RetryCount:  0,
			CreatedAt:   time.Now(),
		}

		if err := s.outboxRepo.CreateEvent(ctx, outboxEvent); err != nil {
			return nil, fmt.Errorf("failed to create outbox event: %w", err)
		}
	}

	return message, nil
}

func (s *service) ListMessages(ctx context.Context, sessionID string) ([]entity.ChatMessage, error) {
	return s.repo.ListMessagesBySession(ctx, sessionID)
}

func (s *service) DeleteMessage(ctx context.Context, id string) error {
	return s.repo.DeleteMessage(ctx, id)
}
