package client

import (
	"context"
	"fmt"
	"portfolio-ai/pkg/config"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Client struct {
	*genai.Client
}

func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	if cfg.AI.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.AI.GeminiAPIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}
	return &Client{client}, nil
}
