package milvus

import (
	"context"
	"fmt"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

const (
    CollectionDanKnowledge       = "dan_knowledge"
    CollectionDanVisitorMemory  = "dan_visitor_memory"
    DimEmbedding                = 3072 // Dimension for Gemini embeddings
)

func (c *Client) InitCollection(ctx context.Context) error {
    if err := c.ensureKnowledgeCollection(ctx); err != nil {
        return err
    }
    if err := c.ensureVisitorMemoryCollection(ctx); err != nil {
        return err
    }
    return nil
}

func (c *Client) ensureKnowledgeCollection(ctx context.Context) error {
    has, err := c.HasCollection(ctx, CollectionDanKnowledge)
    if err != nil {
        return fmt.Errorf("failed to check collection: %w", err)
    }
    if has {
        return nil
    }

    schema := &entity.Schema{
        CollectionName: CollectionDanKnowledge,
        Description:    "Knowledge chunks for dan semantic search",
        AutoID:         false,
        Fields: []*entity.Field{
            {
                Name:       "chunk_id",
                DataType:   entity.FieldTypeVarChar,
                PrimaryKey: true,
                TypeParams: map[string]string{"max_length": "26"},
            },
            {
                Name:       "document_id",
                DataType:   entity.FieldTypeVarChar,
                TypeParams: map[string]string{"max_length": "26"},
            },
            {
                Name:       "source_type",
                DataType:   entity.FieldTypeVarChar,
                TypeParams: map[string]string{"max_length": "50"},
            },
            {
                Name:       "source_id",
                DataType:   entity.FieldTypeVarChar,
                TypeParams: map[string]string{"max_length": "26"},
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

    idx, err := entity.NewIndexAUTOINDEX(entity.COSINE)
    if err != nil {
        return fmt.Errorf("failed to create index definition: %w", err)
    }

    if err := c.CreateIndex(ctx, CollectionDanKnowledge, "embedding", idx, false); err != nil {
        return fmt.Errorf("failed to create index: %w", err)
    }

    if err := c.LoadCollection(ctx, CollectionDanKnowledge, false); err != nil {
        return fmt.Errorf("failed to load collection: %w", err)
    }

    return nil
}

func (c *Client) ensureVisitorMemoryCollection(ctx context.Context) error {
    has, err := c.HasCollection(ctx, CollectionDanVisitorMemory)
    if err != nil {
        return fmt.Errorf("failed to check visitor memory collection: %w", err)
    }
    if has {
        return nil
    }

    schema := &entity.Schema{
        CollectionName: CollectionDanVisitorMemory,
        Description:    "Visitor memories for personalized retrieval",
        AutoID:         false,
        Fields: []*entity.Field{
            {
                Name:       "memory_id",
                DataType:   entity.FieldTypeVarChar,
                PrimaryKey: true,
                TypeParams: map[string]string{"max_length": "26"},
            },
            {
                Name:       "visitor_id",
                DataType:   entity.FieldTypeVarChar,
                TypeParams: map[string]string{"max_length": "26"},
            },
            {
                Name:       "category",
                DataType:   entity.FieldTypeVarChar,
                TypeParams: map[string]string{"max_length": "50"},
            },
            {
                Name:       "key",
                DataType:   entity.FieldTypeVarChar,
                TypeParams: map[string]string{"max_length": "128"},
            },
            {
                Name:       "value",
                DataType:   entity.FieldTypeVarChar,
                TypeParams: map[string]string{"max_length": "1024"},
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
        return fmt.Errorf("failed to create visitor memory collection: %w", err)
    }

    idx, err := entity.NewIndexAUTOINDEX(entity.COSINE)
    if err != nil {
        return fmt.Errorf("failed to create visitor memory index definition: %w", err)
    }

    if err := c.CreateIndex(ctx, CollectionDanVisitorMemory, "embedding", idx, false); err != nil {
        return fmt.Errorf("failed to create visitor memory index: %w", err)
    }

    if err := c.LoadCollection(ctx, CollectionDanVisitorMemory, false); err != nil {
        return fmt.Errorf("failed to load visitor memory collection: %w", err)
    }

    return nil
}
