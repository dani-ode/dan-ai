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

	_, err := c.Upsert(ctx, CollectionPortfolioKnowledge, "", chunkIDCol, docIDCol, sourceTypeCol, sourceIDCol, embedCol)
	if err != nil {
		return fmt.Errorf("failed to upsert vectors to milvus: %w", err)
	}

	return nil
}

func (c *Client) DeleteVectorsByDocumentID(ctx context.Context, documentID string) error {
	expr := fmt.Sprintf(`document_id == "%s"`, documentID)
	err := c.Delete(ctx, CollectionPortfolioKnowledge, "", expr)
	if err != nil {
		return fmt.Errorf("failed to delete vectors from milvus: %w", err)
	}
	return nil
}
