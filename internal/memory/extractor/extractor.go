package extractor

import (
	"context"
	"fmt"
	"strings"

	chatEntity "dan-ai/internal/chat/entity"
	memoryEntity "dan-ai/internal/memory/entity"
)

type Extractor interface {
	ExtractMemories(ctx context.Context, visitorID string, messages []chatEntity.ChatMessage, assistantMessageID string) ([]memoryEntity.Memory, error)
}

type extractor struct{}

func NewExtractor() Extractor {
	return &extractor{}
}

func (e *extractor) ExtractMemories(ctx context.Context, visitorID string, messages []chatEntity.ChatMessage, assistantMessageID string) ([]memoryEntity.Memory, error) {
	if len(messages) == 0 {
		return nil, nil
	}

	var assistantMessage *chatEntity.ChatMessage
	for _, msg := range messages {
		if msg.ID == assistantMessageID && msg.Role == "assistant" {
			assistantMessage = &msg
			break
		}
	}

	if assistantMessage == nil {
		return nil, fmt.Errorf("assistant message %s not found", assistantMessageID)
	}

	recentMessages := make([]string, 0, len(messages))
	for _, msg := range messages {
		recentMessages = append(recentMessages, fmt.Sprintf("%s: %s", msg.Role, strings.TrimSpace(msg.Content)))
	}

	summary := strings.Join(recentMessages, "\n")

	memories := []memoryEntity.Memory{
		{
			VisitorID:  visitorID,
			Category:   "chat",
			Key:        "recent_conversation",
			Value:      summary,
			Confidence: 0.75,
		},
		{
			VisitorID:  visitorID,
			Category:   "assistant_response",
			Key:        fmt.Sprintf("assistant_response_%s", assistantMessage.ID),
			Value:      assistantMessage.Content,
			Confidence: 0.75,
		},
	}

	return memories, nil
}
