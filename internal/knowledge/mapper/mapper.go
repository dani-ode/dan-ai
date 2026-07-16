package mapper

import (
	"dan-ai/internal/knowledge/entity"
	pb "dan-ai/proto/knowledge"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func DocumentEntityToProto(doc *entity.KnowledgeDocument) *pb.KnowledgeDocument {
	if doc == nil {
		return nil
	}

	p := &pb.KnowledgeDocument{
		Id:             doc.ID,
		SourceType:     doc.SourceType,
		SourceId:       doc.SourceID,
		Title:          doc.Title,
		Content:        doc.Content,
		Checksum:       doc.Checksum,
		Version:        doc.Version,
		Status:         doc.Status,
		EmbeddingModel: doc.EmbeddingModel,
		CreatedAt:      timestamppb.New(doc.CreatedAt),
		UpdatedAt:      timestamppb.New(doc.UpdatedAt),
	}

	if doc.LastEmbeddedAt != nil {
		p.LastEmbeddedAt = timestamppb.New(*doc.LastEmbeddedAt)
	}

	return p
}

func ChunkEntityToProto(chunk *entity.KnowledgeChunk) *pb.KnowledgeChunk {
	if chunk == nil {
		return nil
	}

	return &pb.KnowledgeChunk{
		Id:             chunk.ID,
		DocumentId:     chunk.DocumentID,
		ChunkIndex:     chunk.ChunkIndex,
		Content:        chunk.Content,
		TokenCount:     chunk.TokenCount,
		EmbeddingModel: chunk.EmbeddingModel,
		CreatedAt:      timestamppb.New(chunk.CreatedAt),
	}
}
