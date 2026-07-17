package provider

import (
	"context"
	"dan-ai/internal/ai/schema"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type openaiProvider struct {
	client openai.Client
}

func NewOpenAIProvider(apiKey string) Provider {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &openaiProvider{
		client: client,
	}
}

func (p *openaiProvider) GenerateChunks(ctx context.Context, modelName, systemInstruction, documentTitle, documentContent string) ([]schema.Chunk, error) {
	prompt := fmt.Sprintf("Document Title: %s\n\nDocument Content:\n%s", documentTitle, documentContent)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemInstruction),
		openai.UserMessage(prompt),
	}

	resp, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    modelName,
		Messages: messages,
	})
	if err != nil {
		return nil, fmt.Errorf("openai chunk generation failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no content returned from openai")
	}

	rawText := resp.Choices[0].Message.Content

	// Robust JSON extraction
	startIdx := strings.Index(rawText, "{")
	endIdx := strings.LastIndex(rawText, "}")
	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		rawText = rawText[startIdx : endIdx+1]
	} else {
		startIdx = strings.Index(rawText, "[")
		endIdx = strings.LastIndex(rawText, "]")
		if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
			rawText = rawText[startIdx : endIdx+1]
		}
	}

	var result schema.KnowledgeBuilderResponse
	if err := json.Unmarshal([]byte(rawText), &result); err != nil {
		log.Printf("Failed raw JSON: %s", rawText)
		return nil, fmt.Errorf("failed to parse openai json: %w", err)
	}

	return result.Chunks, nil
}

func (p *openaiProvider) GenerateEmbedding(ctx context.Context, modelName, text string) ([]float32, error) {
	resp, err := p.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: modelName,
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(text),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("openai embedding failed: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned from openai")
	}

	// Convert float64 to float32
	embedding64 := resp.Data[0].Embedding
	embedding := make([]float32, len(embedding64))
	for i, v := range embedding64 {
		embedding[i] = float32(v)
	}

	return embedding, nil
}

func (p *openaiProvider) GenerateChatResponse(ctx context.Context, modelName, systemInstruction string, prompt string) (*ChatResponse, error) {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(prompt),
	}

	if systemInstruction != "" {
		messages = []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemInstruction),
			openai.UserMessage(prompt),
		}
	}

	resp, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    modelName,
		Messages: messages,
	})
	if err != nil {
		return nil, fmt.Errorf("openai generate content failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no content returned from openai")
	}

	var promptTokens, completionTokens int32
	if resp.Usage.PromptTokens > 0 {
		promptTokens = int32(resp.Usage.PromptTokens)
	}
	if resp.Usage.CompletionTokens > 0 {
		completionTokens = int32(resp.Usage.CompletionTokens)
	}

	return &ChatResponse{
		Content:          resp.Choices[0].Message.Content,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
	}, nil
}
