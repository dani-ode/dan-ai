package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"portfolio-ai/internal/ai/client"
	"portfolio-ai/internal/ai/schema"

	"github.com/google/generative-ai-go/genai"
)

type Provider interface {
	GenerateChunks(ctx context.Context, documentTitle, documentContent string) ([]schema.Chunk, error)
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
}

type geminiProvider struct {
	client *client.Client
}

func NewGeminiProvider(client *client.Client) Provider {
	return &geminiProvider{
		client: client,
	}
}

func (p *geminiProvider) GenerateChunks(ctx context.Context, documentTitle, documentContent string) ([]schema.Chunk, error) {
	model := p.client.GenerativeModel("gemini-3.5-flash")
	model.ResponseMIMEType = "application/json"

	// Construct JSON schema manually or rely on basic JSON mode.
	// Gemini JSON schema syntax requires structured format, but for simplicity we rely on a strong system prompt + MIME type
	model.SystemInstruction = genai.NewUserContent(genai.Text(`You are an expert knowledge extractor. 
Your task is to chunk the provided document into self-contained segments suitable for vector search.
Return a JSON object with a "chunks" array. Each chunk must have "title", "content" (the detailed text), and "keywords" (array of strings).`))

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

func (p *geminiProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	em := p.client.EmbeddingModel("gemini-embedding-2")
	res, err := em.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		return nil, fmt.Errorf("gemini embedding failed: %w", err)
	}

	if res.Embedding == nil || len(res.Embedding.Values) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return res.Embedding.Values, nil
}
