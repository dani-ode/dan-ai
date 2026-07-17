package provider

import (
	"context"
	"dan-ai/internal/ai/client"
	"dan-ai/internal/ai/schema"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

type ChatResponse struct {
	Content          string
	PromptTokens     int32
	CompletionTokens int32
}

type Provider interface {
	GenerateChunks(ctx context.Context, modelName, systemInstruction, documentTitle, documentContent string) ([]schema.Chunk, error)
	GenerateEmbedding(ctx context.Context, modelName, text string) ([]float32, error)
	GenerateChatResponse(ctx context.Context, modelName, systemInstruction string, prompt string) (*ChatResponse, error)
}

type geminiProvider struct {
	client *client.Client
}

func NewGeminiProvider(client *client.Client) Provider {
	return &geminiProvider{
		client: client,
	}
}

func (p *geminiProvider) GenerateChunks(ctx context.Context, modelName, systemInstruction, documentTitle, documentContent string) ([]schema.Chunk, error) {
	model := p.client.GenerativeModel(modelName)
	model.ResponseMIMEType = "application/json"

	// Use system instruction from the prompt config
	if systemInstruction != "" {
		model.SystemInstruction = genai.NewUserContent(genai.Text(systemInstruction))
	}

	prompt := fmt.Sprintf("Document Title: %s\n\nDocument Content:\n%s", documentTitle, documentContent)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("gemini chunk generation failed: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content returned from gemini")
	}

	part := resp.Candidates[0].Content.Parts[0]
	textPart, ok := part.(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unexpected response type from gemini")
	}

	rawText := string(textPart)

	// Robust JSON extraction: find first '{' and last '}'
	startIdx := strings.Index(rawText, "{")
	endIdx := strings.LastIndex(rawText, "}")
	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		rawText = rawText[startIdx : endIdx+1]
	} else {
		// Fallback if not an object, maybe it's an array?
		startIdx = strings.Index(rawText, "[")
		endIdx = strings.LastIndex(rawText, "]")
		if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
			rawText = rawText[startIdx : endIdx+1]
		}
	}

	var result schema.KnowledgeBuilderResponse
	if err := json.Unmarshal([]byte(rawText), &result); err != nil {
		log.Printf("Failed raw JSON: %s", rawText)
		return nil, fmt.Errorf("failed to parse gemini json: %w", err)
	}

	return result.Chunks, nil
}

func (p *geminiProvider) GenerateEmbedding(ctx context.Context, modelName, text string) ([]float32, error) {
	em := p.client.EmbeddingModel(modelName)
	res, err := em.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		return nil, fmt.Errorf("gemini embedding failed: %w", err)
	}

	if res.Embedding == nil || len(res.Embedding.Values) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return res.Embedding.Values, nil
}

func (p *geminiProvider) GenerateChatResponse(ctx context.Context, modelName, systemInstruction string, prompt string) (*ChatResponse, error) {
	model := p.client.GenerativeModel(modelName)
	if systemInstruction != "" {
		model.SystemInstruction = genai.NewUserContent(genai.Text(systemInstruction))
	}

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("gemini generate content failed: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content returned from gemini")
	}

	part := resp.Candidates[0].Content.Parts[0]
	textPart, ok := part.(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unexpected response type from gemini")
	}

	var promptTokens, completionTokens int32
	if resp.UsageMetadata != nil {
		promptTokens = resp.UsageMetadata.PromptTokenCount
		completionTokens = resp.UsageMetadata.CandidatesTokenCount
	}

	return &ChatResponse{
		Content:          string(textPart),
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
	}, nil
}
