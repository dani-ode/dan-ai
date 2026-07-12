package milvus

import (
	"context"
	"fmt"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

const (
	CollectionPortfolioKnowledge = "portfolio_knowledge"
	DimEmbedding                 = 3072 // Dimension for Gemini embeddings
)

func (c *Client) InitCollection(ctx context.Context) error {
	has, err := c.HasCollection(ctx, CollectionPortfolioKnowledge)
	if err != nil {
		return fmt.Errorf("failed to check collection: %w", err)
	}

	if has {
		return nil // Already exists
	}

	schema := &entity.Schema{
		CollectionName: CollectionPortfolioKnowledge,
		Description:    "Knowledge chunks for portfolio semantic search",
		AutoID:         false,
		Fields: []*entity.Field{
			{
				Name:       "chunk_id",
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: true,
				TypeParams: map[string]string{"max_length": "26"}, // ULID
			},
			{
				Name:       "document_id",
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "26"}, // ULID
			},
			{
				Name:       "source_type",
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "50"},
			},
			{
				Name:       "source_id",
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "26"}, // ULID
			},
			{
				Name:     "embedding",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": fmt.Sprintf("%d", DimEmbedding),
				},
			},
		},
	}

	if err := c.CreateCollection(ctx, schema, entity.DefaultShardNumber); err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	// Create Index for embedding field
	idx, err := entity.NewIndexAUTOINDEX(entity.COSINE)
	if err != nil {
		return fmt.Errorf("failed to create index definition: %w", err)
	}

	if err := c.CreateIndex(ctx, CollectionPortfolioKnowledge, "embedding", idx, false); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// Load the collection into memory
	if err := c.LoadCollection(ctx, CollectionPortfolioKnowledge, false); err != nil {
		return fmt.Errorf("failed to load collection: %w", err)
	}

	return nil
}
