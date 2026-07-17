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
	Score     float32
}

func (c *Client) UpsertVectors(ctx context.Context, collectionName string, vectors []KnowledgeVector) error {
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
	embedCol := entity.NewColumnFloatVector("embedding", len(embeddings[0]), embeddings)

	_, err := c.Upsert(ctx, collectionName, "", chunkIDCol, docIDCol, sourceTypeCol, sourceIDCol, embedCol)
	if err != nil {
		return fmt.Errorf("failed to upsert vectors to milvus: %w", err)
	}

	return nil
}

func (c *Client) DeleteVectorsByDocumentID(ctx context.Context, collectionName string, documentID string) error {
	expr := fmt.Sprintf(`document_id == "%s"`, documentID)
	err := c.Delete(ctx, collectionName, "", expr)
	if err != nil {
		return fmt.Errorf("failed to delete vectors from milvus: %w", err)
	}
	return nil
}

// UpsertVisitorMemoryVectors upserts visitor memory vectors into the visitor memory collection.
func (c *Client) UpsertVisitorMemoryVectors(ctx context.Context, collectionName string, vectors []VisitorMemoryVector) error {
	if len(vectors) == 0 {
		return nil
	}

	memIDs := make([]string, 0, len(vectors))
	visitorIDs := make([]string, 0, len(vectors))
	embeddings := make([][]float32, 0, len(vectors))

	for _, v := range vectors {
		memIDs = append(memIDs, v.MemoryID)
		visitorIDs = append(visitorIDs, v.VisitorID)
		embeddings = append(embeddings, v.Embedding)
	}

	memIDCol := entity.NewColumnVarChar("id", memIDs)
	visitorIDCol := entity.NewColumnVarChar("visitor_id", visitorIDs)
	embedCol := entity.NewColumnFloatVector("embedding", len(embeddings[0]), embeddings)

	_, err := c.Upsert(ctx, collectionName, "", memIDCol, visitorIDCol, embedCol)
	if err != nil {
		return fmt.Errorf("failed to upsert visitor memory vectors to milvus: %w", err)
	}
	return nil
}

// DeleteVisitorMemoryByVisitorID deletes all visitor memory vectors for a given visitor.
func (c *Client) DeleteVisitorMemoryByVisitorID(ctx context.Context, collectionName string, visitorID string) error {
	expr := fmt.Sprintf(`visitor_id == "%s"`, visitorID)
	if err := c.Delete(ctx, collectionName, "", expr); err != nil {
		return fmt.Errorf("failed to delete visitor memory vectors from milvus: %w", err)
	}
	return nil
}

// SearchKnowledge searches the dan_knowledge collection.
func (c *Client) SearchKnowledge(ctx context.Context, collectionName string, queryVector []float32, topK int) ([]KnowledgeVector, error) {
	sp, err := entity.NewIndexAUTOINDEXSearchParam(1)
	if err != nil {
		return nil, fmt.Errorf("failed to create search param: %w", err)
	}

	searchResult, err := c.Search(
		ctx,
		collectionName,
		nil, // partitionNames
		"",  // expr
		[]string{"chunk_id", "document_id", "source_type", "source_id"},
		[]entity.Vector{entity.FloatVector(queryVector)},
		"embedding",
		entity.COSINE,
		topK,
		sp,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search knowledge: %w", err)
	}

	if len(searchResult) == 0 {
		return nil, nil
	}

	res := searchResult[0]
	var vectors []KnowledgeVector

	chunkIDCol := res.Fields.GetColumn("chunk_id")
	docIDCol := res.Fields.GetColumn("document_id")
	sourceTypeCol := res.Fields.GetColumn("source_type")
	sourceIDCol := res.Fields.GetColumn("source_id")

	for i := 0; i < res.ResultCount; i++ {
		chunkID, _ := chunkIDCol.GetAsString(i)
		docID, _ := docIDCol.GetAsString(i)
		sourceType, _ := sourceTypeCol.GetAsString(i)
		sourceID, _ := sourceIDCol.GetAsString(i)

		vectors = append(vectors, KnowledgeVector{
			ChunkID:    chunkID,
			DocumentID: docID,
			SourceType: sourceType,
			SourceID:   sourceID,
		})
	}

	return vectors, nil
}

// SearchVisitorMemory searches the visitor_knowledge collection.
func (c *Client) SearchVisitorMemory(ctx context.Context, collectionName string, visitorID string, queryVector []float32, topK int) ([]VisitorMemoryVector, error) {
	sp, err := entity.NewIndexAUTOINDEXSearchParam(1)
	if err != nil {
		return nil, fmt.Errorf("failed to create search param: %w", err)
	}

	expr := fmt.Sprintf(`visitor_id == "%s"`, visitorID)

	searchResult, err := c.Search(
		ctx,
		collectionName,
		nil, // partitionNames
		expr,
		[]string{"id", "visitor_id"},
		[]entity.Vector{entity.FloatVector(queryVector)},
		"embedding",
		entity.COSINE,
		topK,
		sp,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search visitor memory: %w", err)
	}

	if len(searchResult) == 0 {
		return nil, nil
	}

	res := searchResult[0]
	var vectors []VisitorMemoryVector

	memIDCol := res.Fields.GetColumn("id")
	visitorIDCol := res.Fields.GetColumn("visitor_id")

	for i := 0; i < res.ResultCount; i++ {
		memID, _ := memIDCol.GetAsString(i)
		vID, _ := visitorIDCol.GetAsString(i)

		var score float32
		if i < len(res.Scores) {
			score = res.Scores[i]
		}

		vectors = append(vectors, VisitorMemoryVector{
			MemoryID:  memID,
			VisitorID: vID,
			Score:     score,
		})
	}

	return vectors, nil
}
