package milvus

import (
	"context"
	"fmt"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type KnowledgeVector struct {
	ChunkID    string
	DocumentID string
	SourceType string
	SourceID   string
	Embedding  []float32
}

type VisitorMemoryVector struct {
	MemoryID  string
	VisitorID string
	Category  string
	Key       string
	Value     string
	Embedding []float32
}

func (c *Client) UpsertVectors(ctx context.Context, vectors []KnowledgeVector) error {
	if len(vectors) == 0 {
		return nil
	}

	chunkIDs := make([]string, 0, len(vectors))
	documentIDs := make([]string, 0, len(vectors))
	sourceTypes := make([]string, 0, len(vectors))
	sourceIDs := make([]string, 0, len(vectors))
	embeddings := make([][]float32, 0, len(vectors))

	for _, v := range vectors {
		chunkIDs = append(chunkIDs, v.ChunkID)
		documentIDs = append(documentIDs, v.DocumentID)
		sourceTypes = append(sourceTypes, v.SourceType)
		sourceIDs = append(sourceIDs, v.SourceID)
		embeddings = append(embeddings, v.Embedding)
	}

	chunkIDCol := entity.NewColumnVarChar("chunk_id", chunkIDs)
	docIDCol := entity.NewColumnVarChar("document_id", documentIDs)
	sourceTypeCol := entity.NewColumnVarChar("source_type", sourceTypes)
	sourceIDCol := entity.NewColumnVarChar("source_id", sourceIDs)
	embedCol := entity.NewColumnFloatVector("embedding", DimEmbedding, embeddings)

	_, err := c.Upsert(ctx, CollectionDanKnowledge, "", chunkIDCol, docIDCol, sourceTypeCol, sourceIDCol, embedCol)
	if err != nil {
		return fmt.Errorf("failed to upsert vectors to milvus: %w", err)
	}

	return nil
}

func (c *Client) DeleteVectorsByDocumentID(ctx context.Context, documentID string) error {
	expr := fmt.Sprintf(`document_id == "%s"`, documentID)
	err := c.Delete(ctx, CollectionDanKnowledge, "", expr)
	if err != nil {
		return fmt.Errorf("failed to delete vectors from milvus: %w", err)
	}
	return nil
}

// UpsertVisitorMemoryVectors upserts visitor memory vectors into the visitor memory collection.
func (c *Client) UpsertVisitorMemoryVectors(ctx context.Context, vectors []VisitorMemoryVector) error {
	if len(vectors) == 0 {
		return nil
	}

	memIDs := make([]string, 0, len(vectors))
	visitorIDs := make([]string, 0, len(vectors))
	categories := make([]string, 0, len(vectors))
	keys := make([]string, 0, len(vectors))
	values := make([]string, 0, len(vectors))
	embeddings := make([][]float32, 0, len(vectors))

	for _, v := range vectors {
		memIDs = append(memIDs, v.MemoryID)
		visitorIDs = append(visitorIDs, v.VisitorID)
		categories = append(categories, v.Category)
		keys = append(keys, v.Key)
		values = append(values, v.Value)
		embeddings = append(embeddings, v.Embedding)
	}

	memIDCol := entity.NewColumnVarChar("memory_id", memIDs)
	visitorIDCol := entity.NewColumnVarChar("visitor_id", visitorIDs)
	categoryCol := entity.NewColumnVarChar("category", categories)
	keyCol := entity.NewColumnVarChar("key", keys)
	valueCol := entity.NewColumnVarChar("value", values)
	embedCol := entity.NewColumnFloatVector("embedding", DimEmbedding, embeddings)

	_, err := c.Upsert(ctx, CollectionDanVisitorMemory, "", memIDCol, visitorIDCol, categoryCol, keyCol, valueCol, embedCol)
	if err != nil {
		return fmt.Errorf("failed to upsert visitor memory vectors to milvus: %w", err)
	}
	return nil
}

// DeleteVisitorMemoryByVisitorID deletes all visitor memory vectors for a given visitor.
func (c *Client) DeleteVisitorMemoryByVisitorID(ctx context.Context, visitorID string) error {
	expr := fmt.Sprintf(`visitor_id == "%s"`, visitorID)
	if err := c.Delete(ctx, CollectionDanVisitorMemory, "", expr); err != nil {
		return fmt.Errorf("failed to delete visitor memory vectors from milvus: %w", err)
	}
	return nil
}
