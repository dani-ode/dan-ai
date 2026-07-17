package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"dan-ai/internal/chat/repository"
	"dan-ai/internal/memory/extractor"
	"dan-ai/internal/memory/service"
	promptrepo "dan-ai/internal/prompt/repository"
	"dan-ai/pkg/kafka"
)

const defaultMemoryModel = "gemini-3.1-flash-lite"

type Processor struct {
	chatRepo   repository.Repository
	memorySvc  service.Service
	extractor  extractor.Extractor
	promptRepo promptrepo.Repository
}

func NewProcessor(chatRepo repository.Repository, memorySvc service.Service, extractor extractor.Extractor, promptRepo promptrepo.Repository) *Processor {
	return &Processor{
		chatRepo:   chatRepo,
		memorySvc:  memorySvc,
		extractor:  extractor,
		promptRepo: promptRepo,
	}
}

type chatCompletedPayload struct {
	VisitorID          string `json:"visitor_id"`
	SessionID          string `json:"session_id"`
	AssistantMessageID string `json:"assistant_message_id"`
	PromptID           string `json:"prompt_id"`
}

func (p *Processor) ProcessEvent(ctx context.Context, event kafka.Event) error {
	if event.EventType != "chat.completed" {
		return nil
	}

	var payload chatCompletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal chat.completed payload: %w", err)
	}

	// Resolve model name from prompt_id
	modelName := defaultMemoryModel
	if payload.PromptID != "" {
		prompt, err := p.promptRepo.Get(ctx, payload.PromptID)
		if err == nil && prompt.AIModel.Name != "" {
			modelName = prompt.AIModel.Name
		}
	}

	messages, err := p.chatRepo.ListMessagesBySession(ctx, payload.SessionID)
	if err != nil {
		return fmt.Errorf("failed to list messages for session %s: %w", payload.SessionID, err)
	}

	if len(messages) == 0 {
		log.Printf("no messages found for session %s", payload.SessionID)
		return nil
	}

	memories, err := p.extractor.ExtractMemories(ctx, modelName, payload.VisitorID, messages, payload.AssistantMessageID)
	if err != nil {
		return fmt.Errorf("failed to extract memories: %w", err)
	}

	if len(memories) == 0 {
		log.Printf("no memories extracted for session %s", payload.SessionID)
		return nil
	}

	if err := p.memorySvc.SaveMemories(ctx, modelName, memories); err != nil {
		return fmt.Errorf("failed to save memories: %w", err)
	}

	log.Printf("saved %d memory records for visitor %s using model %s", len(memories), payload.VisitorID, modelName)
	return nil
}
