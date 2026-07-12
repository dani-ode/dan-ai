package milvus

import (
	"context"
	"fmt"
	"portfolio-ai/pkg/config"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

type Client struct {
	client.Client
}

func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	c, err := client.NewClient(ctx, client.Config{
		Address: cfg.Milvus.Address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Milvus: %w", err)
	}

	return &Client{c}, nil
}
