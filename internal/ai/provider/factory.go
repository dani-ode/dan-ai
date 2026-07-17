package provider

import (
	"context"
	"fmt"

	"dan-ai/internal/ai/client"
	"dan-ai/pkg/config"
)

// NewProviderByName creates the appropriate AI provider based on the provider name.
// Supported providers: "gemini", "openai".
func NewProviderByName(ctx context.Context, providerName string, cfg *config.Config) (Provider, error) {
	switch providerName {
	case "gemini":
		genaiClient, err := client.NewClient(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create gemini client: %w", err)
		}
		return NewGeminiProvider(genaiClient), nil
	case "openai":
		if cfg.AI.OpenAIAPIKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY is not set")
		}
		return NewOpenAIProvider(cfg.AI.OpenAIAPIKey), nil
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", providerName)
	}
}
