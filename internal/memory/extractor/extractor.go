package extractor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	aiprovider "dan-ai/internal/ai/provider"
	chatEntity "dan-ai/internal/chat/entity"
	memoryEntity "dan-ai/internal/memory/entity"
	promptrepo "dan-ai/internal/prompt/repository"
	"dan-ai/pkg/ulid"
)

type Extractor interface {
	ExtractMemories(ctx context.Context, modelName, visitorID string, messages []chatEntity.ChatMessage, assistantMessageID string) ([]memoryEntity.Memory, error)
}

type extractor struct {
	aiRegistry *aiprovider.Registry
	promptRepo promptrepo.Repository
}

func NewExtractor(aiRegistry *aiprovider.Registry, promptRepo promptrepo.Repository) Extractor {
	return &extractor{
		aiRegistry: aiRegistry,
		promptRepo: promptRepo,
	}
}

type ExtractedMemory struct {
	Save       bool   `json:"save"`
	Importance int    `json:"importance"`
	Category   string `json:"category"`
	Key        string `json:"key"`
	Memory     string `json:"memory"`
}

const ExtractorSystemInstruction = `You are a memory extractor.
Your task is to analyze the conversation between the User and the Assistant and extract any useful facts or preferences about the User that would be valuable for future personalized conversations.

Ignore:
- greetings
- thanks
- jokes
- generic questions

Identify if there is any long-term memory to save.
The key should be a short, unique key/slug in lowercase with hyphens (e.g. "dan-ai-kafka" or "location-surabaya" or "favorite-language-golang").

Return a JSON object in this exact format:
{
  "save": true|false,
  "importance": 1-5,
  "category": "project|experience|certificate|skill|etc",
  "key": "short-unique-key-slug",
  "memory": "concise description of the visitor context (max 250 characters)"
}`

func (e *extractor) ExtractMemories(ctx context.Context, modelName, visitorID string, messages []chatEntity.ChatMessage, assistantMessageID string) ([]memoryEntity.Memory, error) {
	if len(messages) == 0 {
		return nil, nil
	}

	// Resolve dynamic system prompt from database
	systemInstruction := ExtractorSystemInstruction
	allPrompts, err := e.promptRepo.List(ctx, false)
	if err == nil {
		for _, p := range allPrompts {
			if p.Name == "Memory Extractor" {
				systemInstruction = p.SystemPrompt
				break
			}
		}
	}

	var dialogueLines []string
	for _, msg := range messages {
		roleName := "User"
		if msg.Role == "assistant" {
			roleName = "Assistant"
		}
		dialogueLines = append(dialogueLines, fmt.Sprintf("%s: %s", roleName, strings.TrimSpace(msg.Content)))
		if msg.ID == assistantMessageID {
			break
		}
	}
	dialogueText := strings.Join(dialogueLines, "\n")

	prompt := fmt.Sprintf("Dialogue:\n%s", dialogueText)

	// Resolve provider by looking up the prompt's model provider
	providerName := "gemini" // default
	allPrompts2, err := e.promptRepo.List(ctx, false)
	if err == nil {
		for _, p := range allPrompts2 {
			if p.Name == "Memory Extractor" && p.AIModel.Provider != "" {
				providerName = p.AIModel.Provider
				break
			}
		}
	}
	extractProvider, err := e.aiRegistry.Get(providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider %q for extraction: %w", providerName, err)
	}

	chatResp, err := extractProvider.GenerateChatResponse(ctx, modelName, systemInstruction, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to call gemini: %w", err)
	}

	rawText := strings.TrimSpace(chatResp.Content)
	startIdx := strings.Index(rawText, "{")
	endIdx := strings.LastIndex(rawText, "}")
	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		rawText = rawText[startIdx : endIdx+1]
	}

	var extMem ExtractedMemory
	if err := json.Unmarshal([]byte(rawText), &extMem); err != nil {
		return nil, fmt.Errorf("failed to parse extracted memory JSON: %w", err)
	}

	if !extMem.Save || extMem.Memory == "" {
		return nil, nil
	}

	mem := memoryEntity.Memory{
		ID:         ulid.New(),
		VisitorID:  visitorID,
		Category:   extMem.Category,
		MemoryText: extMem.Memory,
		Importance: extMem.Importance,
	}

	return []memoryEntity.Memory{mem}, nil
}
